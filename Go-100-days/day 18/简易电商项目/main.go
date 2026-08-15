// main.go
// 程序入口文件
// 学习要点：Go 程序入口在 main 包的 main 函数
//
// 启动流程：
// 1. 加载配置文件（Viper 读取 config.yaml）
// 2. 初始化 MySQL 连接（GORM）
// 3. 初始化 Redis 连接
// 4. 自动建表（GORM AutoMigrate）
// 5. 注册路由
// 6. 启动 HTTP 服务
//
// 学习要点：Go 的初始化顺序
// - main 函数是入口
// - init 函数在 main 之前执行（本项目没用）
// - 全局变量在 init 之前初始化

package main

import (
	"fmt"
	"log"

	"simple-ecommerce/internal/model"
	"simple-ecommerce/internal/router"
	"simple-ecommerce/initialize"
	"simple-ecommerce/pkg/config"
)

func main() {
	// 1. 加载配置文件
	// 学习要点：config.Load("./config") 在 config 目录下找 config.yaml
	log.Println("开始加载配置...")
	if err := config.Load("./config"); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	log.Printf("配置加载成功，应用: %s", config.GlobalConfig.App.Name)

	// 2. 初始化 MySQL
	if err := initialize.InitMySQL(); err != nil {
		log.Fatalf("初始化 MySQL 失败: %v", err)
	}

	// 3. 初始化 Redis
	if err := initialize.InitRedis(); err != nil {
		log.Fatalf("初始化 Redis 失败: %v", err)
	}

	// 4. 自动建表
	// 学习要点：AutoMigrate 根据模型自动创建/更新表结构
	// 开发环境方便，生产环境建议用 SQL 脚本管理
	log.Println("开始自动建表...")
	if err := initialize.DB.AutoMigrate(
		&model.User{},
		&model.Product{},
		&model.Order{},
		&model.OrderItem{},
	); err != nil {
		log.Fatalf("自动建表失败: %v", err)
	}
	log.Println("建表完成")

	// 5. 注册路由
	r := router.SetupRouter()

	// 6. 启动服务
	port := config.GlobalConfig.App.Port
	addr := fmt.Sprintf(":%d", port)
	log.Printf("服务启动，监听地址: %s", addr)
	log.Printf("访问 http://localhost:%d/api/v1/products 查看商品列表", port)

	// 学习要点：r.Run 启动 HTTP 服务，阻塞主 goroutine
	// 相当于 http.ListenAndServe(addr, r)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
