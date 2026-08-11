package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"blog/internal/handler" // 导入 handler 包
	"blog/internal/middleware"
	"blog/model"      // 导入 model 包
	"blog/pkg/cache"  // 导入 Redis
	"blog/pkg/logger" // 导入 Zap
)

func main() {

	//---------- 0. 初始化日志 ----------
	if err := logger.Init(); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer logger.Sync() // 程序退出前刷新日志缓冲
	//-----------------------------------

	//---------- 1. 加载配置文件 ----------
	viper.SetConfigName("config") // 配置文件名（不用写扩展名）
	viper.SetConfigType("yaml")   // 配置文件类型
	viper.AddConfigPath(".")      // 在当前目录下查找

	if err := viper.ReadInConfig(); err != nil {
		logger.Log.Fatal("加载配置文件失败", zap.String("error", err.Error()))
	}
	//----------------------------------------

	//---------- 2. 连接数据库 ----------
	//从配置里取 database.dsn 的值
	dsn := viper.GetString("database.dsn")

	//用 GORM 打开数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Log.Fatal("连接数据库失败", zap.String("error", err.Error())) // 如果连不上 Mysql, 程序直接退出
	}
	logger.Log.Info("✅ 数据库连接成功")

	// ---------- 2.5 连接 Redis ----------
	redisAddr := viper.GetString("redis.addr")
	if redisAddr == "" {
		redisAddr = "localhost:6379" //默认值
	}

	if err := cache.InitRedis(redisAddr, "", 0); err != nil {
		logger.Log.Fatal("连接 Redis 失败", zap.String("error", err.Error()))
	}
	logger.Log.Info("✅ Redis 连接成功")

	//---------- 3. 自动迁移（移表） ----------
	//检查数据库里有没有 users 表，没有就创建；有就检查字段是否匹配，不匹配就更新
	if err := db.AutoMigrate(&model.User{}, &model.Article{}); err != nil {
		log.Fatalf("数据库迁移失败： %v", err)
	}
	logger.Log.Info("✅ 数据库迁移完成（users 表已就绪）") // ✅ 改为了 Println
	//----------------------------------------
	//---------- 4. 创建用户处理器 ----------
	userhandler := handler.NewUserHandler(db)

	// 5. 创建一个 Gin 路由器 和 一个 handlers
	r := gin.Default()
	articleHandler := handler.NewArticleHandler(db)

	// 6. 注册路由
	//公共路由（不需要登录）
	r.POST("/api/v1/users/register", userhandler.Register)
	r.POST("/api/v1/users/login", userhandler.Login)

	// 文章公开接口
	r.GET("/api/v1/articles", articleHandler.List)       //文章列表
	r.GET("/api/v1/articles/:id", articleHandler.Detail) //文章详情

	//健康检查(公开)
	//定义一个路由：当用户访问 /ping 时，返回 JSON 数据
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"msg": "pong"})
	})

	//-------- 需要登录的路由组 --------
	authGroup := r.Group("api/v1")
	authGroup.Use(middleware.Auth()) //	应用 JWT 认证中间件
	{
		//	文章相关的接口（需要登录）
		authGroup.POST("/articles", articleHandler.Create)       //创建文章
		authGroup.PUT("/articles/:id", articleHandler.Update)    //更新文章
		authGroup.DELETE("/articles/:id", articleHandler.Delete) //删除文章
		// 后续可以继续添加：GET /articles, PUT /articles/:id, DELETE /articles/:id 等

	}

	// 7. 启动服务器
	//从配置中读取端口，如果没有则使用默认值 “:8080”
	port := viper.GetString("server.port")
	if port == "" {
		port = ":8080"
	}
	log.Printf("✅ 服务器启动在 http://localhost%s\n", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}

}
