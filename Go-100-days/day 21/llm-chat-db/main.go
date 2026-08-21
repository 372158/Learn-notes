package main

import (
	"github.com/gin-gonic/gin"

	"llm-chat-db/config"
	"llm-chat-db/db"
	"llm-chat-db/handler"
	"llm-chat-db/middleware"
	"llm-chat-db/service"
)

func main() {
	cfg := config.Load()

	database := db.Init(cfg.DBDSN)

	svc := service.NewChatService(database, cfg.APIKey, cfg.APIURL, cfg.Model)

	// gin.New() 代替 gin.Default()：Default 自带一条默认访问日志，
	// 这里用自定义中间件，日志更简洁、可自己控制格式
	r := gin.New()
	r.Use(gin.Recovery()) // 防止单个 handler panic 拖垮整个服务
	r.Use(middleware.Logger())

	r.GET("/chat", handler.NewChatHandler(svc).Chat)
	r.Run(":" + cfg.Port)
}