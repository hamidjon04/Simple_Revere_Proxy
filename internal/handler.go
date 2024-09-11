package internal

import (
	"fmt"
	"log/slog"
	"proxy/caching"
	"proxy/client"
	"proxy/model"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func ProxyHandler(c *gin.Context, config *model.ConfigSettings, Rdb *redis.Client, logger *slog.Logger) (interface{}, error) {
	requestURL := config.Servers[0] + c.Request.RequestURI
	// Serverlarni ko'paytirib tezlikni va yuklamani yanada oshirish mumkin
	redisDB := caching.NewRedisRepo(Rdb)

	resp, err := redisDB.CheackCache(c, requestURL)
	if err == redis.Nil {
		logger.Info(fmt.Sprintf("%s request uchun response cachingda mavjud emas. So'rov asosiy serverga yuborildi.", requestURL))
		resp, err := client.ForwardBackend(requestURL)
		if err != nil{
			logger.Error(fmt.Sprintf("Server bilan bog'lanishda yoki ma'lumotlarni olishda xatolik: %v", err))
			return nil, err
		}
		return resp, nil
	} else if err != nil {
		logger.Error(fmt.Sprintf("Cachingdan ma'lumotlarni olishda xatolik: %v", err))
		resp, err := client.ForwardBackend(requestURL)
		if err != nil{
			logger.Error(fmt.Sprintf("Server bilan bog'lanishda yoki ma'lumotlarni olishda xatolik: %v", err))
			return nil, err
		}
		return resp, nil
	}
	return resp, err
}
