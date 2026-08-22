package cache

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

// New 连接 Redis，返回客户端；连接失败直接退出（和 db.Init 保持一致风格）
func New(addr string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("连接 Redis 失败: %v", err)
	}
	log.Println("Redis 连接成功:", addr)
	return rdb
}
