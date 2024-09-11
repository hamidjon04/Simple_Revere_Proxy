package main

import (
	"log"
	"proxy/caching"
	"proxy/config"
	"proxy/internal"
	"proxy/logs"

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
		internal.ProxyHandler(ctx, config, Rdb, logger)
	})

	err = r.Run(":8080")
	if err != nil {
		log.Printf("Proxy run bo'lmadi: %v", err)
		panic(err)
	}
}
