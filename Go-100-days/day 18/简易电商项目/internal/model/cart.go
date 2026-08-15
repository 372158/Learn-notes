// internal/model/cart.go
// 购物车数据模型
// 学习要点：购物车用 Redis Hash 存储，这里只定义返回给前端的展示结构
//
// 为什么购物车用 Redis 而非 MySQL：
// 1. 高频读写：加购、改数量频繁，Redis 内存读写极快
// 2. 临时性：购物车数据允许丢失（用户没下单前都是临时的）
// 3. 减轻 MySQL 压力：不占用关系数据库资源
// 4. 天然按用户隔离：用 key cart:{user_id} 天然分用户
//
// Redis Hash 存储结构：
// Key:   cart:{user_id}        例如 cart:1
// Field: {product_id}          例如 100
// Value: {quantity}            例如 2
// 完整示例：cart:1 -> {100:2, 101:1, 102:5}

package model

// CartItem 购物车项（展示用，包含商品详情）
// 学习要点：这个结构体不对应数据库表，仅用于业务传输
type CartItem struct {
	ProductID   uint    `json:"product_id"`   // 商品 ID
	Name        string  `json:"name"`         // 商品名称
	Price       float64 `json:"price"`        // 商品单价
	Quantity    int     `json:"quantity"`     // 购买数量
	ImageURL    string  `json:"image_url"`    // 商品图片
	TotalPrice  float64 `json:"total_price"`  // 小计金额 = price * quantity
}
