// initialize/mysql.go
// MySQL 数据库初始化，使用 GORM 连接 MySQL
// 学习要点：GORM 是 Go 中最流行的 ORM 框架，能像操作对象一样操作数据库
//
// 为什么用 GORM 而非原生 database/sql：
// 1. 自动建表（AutoMigrate）
// 2. 链式调用 API 更优雅
// 3. 自动处理空值、时间等类型转换
// 4. 支持关联查询、事务、钩子等高级功能

package initialize

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"simple-ecommerce/pkg/config"
)

// DB 全局数据库实例
// 学习要点：用全局变量保存数据库连接，整个项目复用同一个连接池
var DB *gorm.DB

// InitMySQL 初始化 MySQL 连接
func InitMySQL() error {
	cfg := config.GlobalConfig.MySQL

	// 学习要点：gorm.Open 接收两个参数
	// 1. mysql.Open(dsn): MySQL 驱动，传入连接字符串
	// 2. gorm.Config: GORM 配置，这里主要配置日志级别
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 日志输出到控制台
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值（超过 1 秒的 SQL 会被标记）
			LogLevel:      logger.Info, // 日志级别：Info 会打印所有执行的 SQL
			Colorful:      true,        // 彩色日志，方便区分
		},
	)

	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return fmt.Errorf("连接 MySQL 失败: %w", err)
	}

	// 配置连接池
	// 学习要点：连接池避免频繁创建/销毁连接，提升性能
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取底层 sql.DB 失败: %w", err)
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns) // 最大空闲连接数
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns) // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour)     // 连接最大存活时间

	DB = db
	log.Println("MySQL 连接成功")
	return nil
}
