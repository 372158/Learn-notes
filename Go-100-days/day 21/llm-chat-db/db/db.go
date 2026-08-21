package db

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"llm-chat-db/models"
)

// Init 连接数据库并自动建表，返回 gorm.DB 实例
func Init(dsn string) *gorm.DB {
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	database.AutoMigrate(&models.Message{})
	log.Println("数据库连接成功")
	return database
}