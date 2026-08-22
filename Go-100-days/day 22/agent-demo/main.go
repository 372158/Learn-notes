package main

import (
	"github.com/gin-gonic/gin"

	"agent-demo/config"
)

func main() {
	cfg := config.Load()
	_ = cfg // 骨架阶段暂不使用全部配置，A8 接线

	r := gin.New()
	r.Use(gin.Recovery())

	// 临时保活检查路由，验证骨架可运行
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "agent-demo 骨架已就绪"})
	})

	port := cfg.Port
	if port == "" {
		port = "8091"
	}
	r.Run(":" + port)
}
