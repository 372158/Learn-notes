// initialize/redis.go
// Redis 初始化，使用 go-redis v9 连接 Redis
// 学习要点：go-redis 是 Go 中最流行的 Redis 客户端
//
// 为什么 Redis 用专门的客户端库而非驱动：
// Redis 是内存数据库，通信协议简单，客户端库封装了命令调用方式，
// 让我们能像调用函数一样执行 Redis 命令（rdb.Set / rdb.Get 等）

package initialize

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"

	"simple-ecommerce/pkg/config"
)

// RDB 全局 Redis 客户端实例
// 学习要点：命名约定，Redis 客户端通常命名为 RDB 或 RedisClient
var RDB *redis.Client

// InitRedis 初始化 Redis 连接
func InitRedis() error {
	cfg := config.GlobalConfig.Redis

	// 学习要点：redis.NewClient 创建客户端，传入配置
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,     // Redis 服务器地址
		Password: cfg.Password, // 密码，没设置则为空
		DB:       cfg.DB,       // 使用的数据库编号（0-15）
	})

	// 学习要点：Ping 测试连接是否正常
	// go-redis v9 的所有命令都需要传 context.Context 参数
	// context.Background() 表示根上下文，没有超时和取消
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("连接 Redis 失败: %w", err)
	}

	RDB = rdb
	log.Println("Redis 连接成功")
	return nil
}
