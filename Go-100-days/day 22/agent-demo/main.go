package main

import (
	"time"

	"github.com/gin-gonic/gin"

	"agent-demo/agent"
	"agent-demo/config"
	"agent-demo/handler"
	"agent-demo/llm"
	"agent-demo/middleware"
	"agent-demo/tools"
)

func main() {
	cfg := config.Load()

	m := llm.NewGLM(cfg.APIKey, cfg.APIURL, cfg.Model)
	mem := agent.NewMemMemory()
	ag := agent.New(m, []tools.Tool{tools.Calc, tools.Weather}, mem, 5)

	secret := cfg.JWTSecret
	if secret == "" { // 没配密钥就用默认，仅限教学演示
		secret = "dev-secret"
	}
	h := handler.New(ag, secret, 24*time.Hour)

	r := gin.New()
	r.Use(gin.Recovery())

	// 无需登录：登录换取 token
	r.POST("/login", h.Login)

	// 受保护：必须带有效 token 才能访问
	authGroup := r.Group("/api")
	authGroup.Use(middleware.Auth(secret))
	{
		authGroup.POST("/chat", h.Chat)
	}

	port := cfg.Port
	if port == "" {
		port = "8091"
	}
	r.Run(":" + port)
}