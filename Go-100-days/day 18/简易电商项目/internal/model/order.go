// internal/model/order.go
// 订单数据模型，包含 Order（订单主表）和 OrderItem（订单项明细）
// 学习要点：一对多关系设计（一个订单包含多个订单项）
//
// 为什么订单要拆成两张表：
// 1. 一个订单可含多个商品，拆表符合数据库范式
// 2. 订单主表存订单级信息（总金额、状态、用户）
// 3. 订单项表存商品级信息（商品快照、单价、数量）
//
// 为什么订单项要存商品快照（名称、单价）：
// 商品价格可能变动，下单时的价格要锁定，不能因商品改价影响历史订单

package model

import "time"

// Order 订单主表
type Order struct {
	ID         uint      `json:"id" gorm:"primaryKey"`                       // 主键
	OrderNo    string    `json:"order_no" gorm:"uniqueIndex;size:32;not null"` // 订单号（业务唯一）
	UserID     uint      `json:"user_id" gorm:"index;not null"`               // 用户 ID（加索引便于查询）
	TotalPrice float64   `json:"total_price" gorm:"type:decimal(10,2)"`        // 订单总金额
	Status     int       `json:"status" gorm:"default:0"`                      // 订单状态：0待支付 1已支付 2已取消
	CreatedAt  time.Time `json:"created_at"`                                   // 创建时间
	UpdatedAt  time.Time `json:"updated_at"`                                   // 更新时间

	// 学习要点：一对多关联
	// - has many: Order 拥有多个 OrderItem
	// - foreignKey: 指定外键字段（OrderItem.OrderID）
	// 查询时用 Preload("Items") 预加载订单项
	Items []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
}

// TableName 自定义表名
func (Order) TableName() string {
	return "orders"
}

// OrderItem 订单项明细表
type OrderItem struct {
	ID         uint    `json:"id" gorm:"primaryKey"`                        // 主键
	OrderID    uint    `json:"order_id" gorm:"index;not null"`               // 所属订单 ID
	ProductID  uint    `json:"product_id" gorm:"not null"`                   // 商品 ID
	ProductName string `json:"product_name" gorm:"size:100;not null"`        // 商品名称快照（下单时锁定）
	Price      float64 `json:"price" gorm:"type:decimal(10,2);not null"`     // 商品单价快照
	Quantity   int     `json:"quantity" gorm:"not null"`                     // 购买数量
	TotalPrice float64 `json:"total_price" gorm:"type:decimal(10,2)"`        // 小计金额
}

// TableName 自定义表名
func (OrderItem) TableName() string {
	return "order_items"
}

// 订单状态常量
// 学习要点：用常量定义状态，避免代码中出现魔法数字
const (
	OrderStatusPending   = 0 // 待支付
	OrderStatusPaid      = 1 // 已支付
	OrderStatusCancelled = 2 // 已取消
)
