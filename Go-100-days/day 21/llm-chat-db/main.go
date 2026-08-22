package main

import (
	"time"

	"github.com/gin-gonic/gin"

	"llm-chat-db/cache"
	"llm-chat-db/config"
	"llm-chat-db/db"
	"llm-chat-db/handler"
	"llm-chat-db/middleware"
	"llm-chat-db/service"
)

func main() {
	cfg := config.Load()

	database := db.Init(cfg.DBDSN)
	rdb := cache.New(cfg.RedisAddr)

	svc := service.NewChatService(database, rdb, cfg.APIKey, cfg.APIURL, cfg.Model)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	// 登录签发 token（不鉴权，任何人都能登录换取自己的 token）
	r.POST("/login", handler.NewLoginHandler(cfg.JWTSecret, 2*time.Hour).Login)

	// /chat 现在受保护：必须先登录拿 token，验证通过才放行
	r.GET("/chat", middleware.Auth(cfg.JWTSecret), handler.NewChatHandler(svc).Chat)

	r.Run(":" + cfg.Port)
}
