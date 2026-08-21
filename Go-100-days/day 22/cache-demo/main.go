package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Message struct {
	ID        uint   `gorm:"primaryKey"`
	SessionID string `gorm:"index"`
	Role      string
	Content   string
}

var ctx = context.Background()

func getMessagesFromDB(db *gorm.DB, sid string) []Message {
	var list []Message
	db.Where("session_id = ?", sid).Find(&list)
	fmt.Println("[DB] 从数据库查了历史") // 标记：真的打到了 DB
	return list
}

func main() {
	// 连 MySQL
	db, err := gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:3306)/chat_app?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Message{})

	// 连 Redis
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()

	// 初始化：往 DB 存几条该会话的记录（模拟已有历史）
	db.Create(&Message{SessionID: "c1", Role: "user", Content: "你好"})
	db.Create(&Message{SessionID: "c1", Role: "assistant", Content: "你好呀"})

	// 查两次，每次都用"先查缓存，miss 再查库"的逻辑
	for i := 0; i < 2; i++ {
		sid := "c1"
		cacheKey := "history:" + sid

		// 空1：先查 Redis 缓存（从缓存读，返回 (val, err)）
		cached, err := rdb.Get(ctx, cacheKey).Result()

		if err == redis.Nil { // 空2：err 等于什么时，表示"缓存未命中"？
			// 未命中 → 查库，回填缓存
			msgs := getMessagesFromDB(db, sid)
			_ = msgs

			// 空3：把查到的历史写进 Redis（用 Set，设个短过期），key 用 cacheKey
			rdb.Set(ctx, cacheKey, "cached-history", 30*time.Second)
			fmt.Println("[缓存] 已回填")
		} else if err != nil {
			panic(err)
		} else {
			// 命中 → 直接用缓存
			fmt.Println("[缓存] 命中:", cached)
		}
	}

	// 空4：查这个 key 的 TTL（看它会不会过期）
	ttl, err := rdb.TTL(ctx, "history:c1").Result()
	// 空5：打印剩余过期时间
	fmt.Println("缓存剩余过期:", ttl)
}