package cache

import "testing"

// TestRedisConnection 验证 Redis 连接 + SET/GET 读写是否正常
func TestRedisConnection(t *testing.T) {
	// 1. 初始化连接（内部会执行 PING 探活）
	if err := InitRedis("localhost:6379", "", 0); err != nil {
		t.Fatalf("❌ Redis 连接失败: %v", err)
	}
	t.Log("✅ Redis PING 成功，连接正常")

	// 2. 写入一个测试键
	key := "blog:test:check"
	if err := Rdb.Set(Ctx, key, "hello-redis", 0).Err(); err != nil {
		t.Fatalf("❌ SET 失败: %v", err)
	}
	t.Log("✅ SET 成功")

	// 3. 读回来比对
	val, err := Rdb.Get(Ctx, key).Result()
	if err != nil {
		t.Fatalf("❌ GET 失败: %v", err)
	}
	if val != "hello-redis" {
		t.Fatalf("❌ 值不匹配: got %q, want %q", val, "hello-redis")
	}
	t.Logf("✅ GET 成功，值为: %s", val)

	// 4. 清理测试键
	Rdb.Del(Ctx, key)
	t.Log("✅ Redis 连接 + 读写全部正常，测试键已清理")
}
