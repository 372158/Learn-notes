package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"shortlink/internal/config"   // 现在会使用
	"shortlink/internal/handler"
	"shortlink/internal/model"
	"shortlink/internal/service"
	"shortlink/pkg/cache"
)

func main() {
	// 1. 加载配置
	if err := config.Load(); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	log.Println("✅ 配置加载成功")

	// 2. 连接 MySQL（使用 config.AppConfig）
	db, err := gorm.Open(mysql.Open(config.AppConfig.Database.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	db.AutoMigrate(&model.ShortLink{})
	log.Println("✅ 数据库连接成功")

	// 3. 连接 Redis
	if err := cache.InitRedis(
		config.AppConfig.Redis.Addr,
		config.AppConfig.Redis.Password,
		config.AppConfig.Redis.DB,
	); err != nil {
		log.Fatalf("连接 Redis 失败: %v", err)
	}

	// 4. 创建 Service 和 Handler
	svc := service.NewShortLinkService(db)
	h := handler.NewShortLinkHandler(svc)

	// 5. 注册路由
	r := gin.Default()
	r.POST("/api/shorten", h.Create)
	r.GET("/:code", h.Redirect)

	// 6. 启动服务
	port := config.AppConfig.Server.Port
	log.Printf("✅ 服务器启动在 http://localhost%s", port)
	r.Run(port)
}