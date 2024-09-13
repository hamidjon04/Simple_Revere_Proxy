package main

import (
	"fmt"
	"log"
	"net/http"
	"proxy/caching"
	"proxy/config"
	"proxy/internal"
	"proxy/logs"
	"proxy/model"

	"github.com/gin-gonic/gin"
)

func main() {
	logger := logs.InitLogger()

	config, err := config.LoadConfig(logger)
	if err != nil {
		log.Printf("Configuratsiya yuklanmadi: %v", err)
		panic(err)
	}

	Rdb := caching.ConnectRedis()
	defer Rdb.Close()

	r := gin.Default()

	r.Any("/*proxyPath", func(ctx *gin.Context) {
		resp, err := internal.ProxyHandler(ctx, config, Rdb, logger)
		if err != nil{
			logger.Error(fmt.Sprintf("Error: %v", err))
			ctx.JSON(http.StatusBadRequest, model.Error{
				Message: "Error!",
			})
			return 
		}
		ctx.JSON(http.StatusOK, resp)
	})

	err = r.Run(":8080")
	if err != nil {
		log.Printf("Proxy run bo'lmadi: %v", err)
		panic(err)
	}
}
