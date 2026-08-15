// internal/dao/cart.go
// 购物车数据访问层（Redis 操作）
// 学习要点：Redis Hash 的常用命令封装
//
// Hash 常用命令对照：
// Redis 命令     | go-redis 方法      | 作用
// HSET          | HSet              | 设置字段值
// HGET          | HGet              | 获取单个字段值
// HGETALL       | HGetAll           | 获取所有字段和值
// HDEL          | HDel              | 删除字段
// HINCRBY       | HIncrBy           | 字段值自增

package dao

import (
	"context"
	"fmt"
	"strconv"

	"simple-ecommerce/initialize"
)

// CartDAO 购物车数据访问对象
type CartDAO struct{}

// NewCartDAO 创建 CartDAO 实例
func NewCartDAO() *CartDAO {
	return &CartDAO{}
}

// cartKey 生成购物车的 Redis key
// 学习要点：用函数封装 key 生成规则，避免拼写错误
// 格式：cart:{user_id}
func (dao *CartDAO) cartKey(userID uint) string {
	return fmt.Sprintf("cart:%d", userID)
}

// AddToCart 添加商品到购物车
// 学习要点：HSet 设置 Hash 字段
// 如果 field 已存在会覆盖，所以加购时要先读再累加（在 service 层处理）
func (dao *CartDAO) AddToCart(userID, productID uint, quantity int) error {
	key := dao.cartKey(userID)
	field := strconv.FormatUint(uint64(productID), 10)
	// HSet 参数：ctx, key, field, value
	return initialize.RDB.HSet(context.Background(), key, field, quantity).Err()
}

// GetQuantity 获取购物车中某商品的数量
// 学习要点：HGet 获取单个字段，返回 *redis.StringCmd
// .Result() 返回 (string, error)，不存在时返回 redis.Nil 错误
func (dao *CartDAO) GetQuantity(userID, productID uint) (int, error) {
	key := dao.cartKey(userID)
	field := strconv.FormatUint(uint64(productID), 10)
	val, err := initialize.RDB.HGet(context.Background(), key, field).Result()
	if err != nil {
		return 0, err // redis.Nil 表示字段不存在
	}
	return strconv.Atoi(val)
}

// GetAllItems 获取购物车所有商品
// 学习要点：HGetAll 返回 map[string]string
// key 是 product_id（字符串），value 是 quantity（字符串）
func (dao *CartDAO) GetAllItems(userID uint) (map[uint]int, error) {
	key := dao.cartKey(userID)
	result, err := initialize.RDB.HGetAll(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	// 转换 map[string]string -> map[uint]int
	items := make(map[uint]int)
	for pidStr, qtyStr := range result {
		pid, err := strconv.ParseUint(pidStr, 10, 64)
		if err != nil {
			continue
		}
		qty, err := strconv.Atoi(qtyStr)
		if err != nil {
			continue
		}
		items[uint(pid)] = qty
	}
	return items, nil
}

// UpdateQuantity 更新购物车商品数量
// 学习要点：和 AddToCart 一样用 HSet，直接覆盖
func (dao *CartDAO) UpdateQuantity(userID, productID uint, quantity int) error {
	key := dao.cartKey(userID)
	field := strconv.FormatUint(uint64(productID), 10)
	return initialize.RDB.HSet(context.Background(), key, field, quantity).Err()
}

// RemoveItem 从购物车移除商品
// 学习要点：HDel 删除 Hash 中的字段
func (dao *CartDAO) RemoveItem(userID, productID uint) error {
	key := dao.cartKey(userID)
	field := strconv.FormatUint(uint64(productID), 10)
	return initialize.RDB.HDel(context.Background(), key, field).Err()
}

// ClearCart 清空购物车
// 学习要点：Del 删除整个 key
func (dao *CartDAO) ClearCart(userID uint) error {
	key := dao.cartKey(userID)
	return initialize.RDB.Del(context.Background(), key).Err()
}
