package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	// ① 创建客户端：地址指向本机 6379（redis-test 容器映射的端口）
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // 无密码， 不填 Password
	})

	defer rdb.Close()


	// ② 验证连接
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic("Redis 连接失败：" + err.Error())
	}
	fmt.Println("Redis 连接成功")

		// Set：写入，带 10 秒过期
	err := rdb.Set(ctx, "greeting", "hello go-redis", 10*time.Second).Err()
	if err != nil {
		panic("Set 失败: " + err.Error())
	}
	fmt.Println("写入成功")

	// Get：读取
	val, err := rdb.Get(ctx, "greeting").Result() // 空1：怎么读出 greeting 的值？（用 Get + Result）
	fmt.Println("读取到:", val)

	// TTL：看还剩多少秒
	ttl, err := rdb.TTL(ctx, "greeting").Result()  // 空2：怎么查过期时间？（用 TTL + Result）
	fmt.Println("剩余过期:", ttl)

	// 等 2 秒再看 TTL，体会"越来越短"
	time.Sleep(2 * time.Second)
	ttl2, err := rdb.TTL(ctx, "greeting").Result()
	fmt.Println("2秒后剩余:", ttl2)

	// 管道批量写（了解即可，性能优化用）
	pipe := rdb.Pipeline()
	pipe.Set(ctx, "a", 1, 0)
	pipe.Set(ctx, "b", 2, 0)
	pipe.Exec(ctx)

	// Hash：存对象（存用户信息）
	rdb.HSet(ctx, "user:1", map[string]interface{}{"name": "小明", "age": 3})
	name, err := rdb.HGet(ctx, "user:1", "name").Result()// 空3：从 Hash 里读 name 字段（用 HGet + Result）
	fmt.Println("user:1 的 name:", name)
}
