package cache

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client
var Ctx = context.Background()

// InitRedis 初始化 Redis 链接
func InitRedis(addr, password string, db int) error {
	Rdb = redis.NewClient(&redis.Options{
		Addr:  addr,		//	例如“localhost :6379"
		Password: password,	//	没有密码就留空 ""
		DB: db,				// 	默认使用 0 号数据库
	})

	//测试链接
	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		return err
	}
	log.Println("✅ Redis 连接成功")
	return nil

}
