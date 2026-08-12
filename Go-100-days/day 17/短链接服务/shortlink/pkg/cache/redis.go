package cache

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"  // 注意是 v9
)

var Rdb *redis.Client
var Ctx = context.Background()

func InitRedis(addr, password string, db int) error {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// v9 的 Ping 方法需要传入 context
	if err := Rdb.Ping(Ctx).Err(); err != nil {
		return err
	}
	log.Println("✅ Redis 连接成功")
	return nil
}