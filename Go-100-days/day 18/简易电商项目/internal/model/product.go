// internal/model/product.go
// 商品数据模型，对应数据库 products 表
// 学习要点：商品模型设计，包含名称、价格、库存等电商核心字段
//
// 字段设计说明：
// - Price 用 float64 简化（生产环境建议用 decimal 类型避免精度问题）
// - Stock 库存，下单时扣减
// - ImageURL 商品图片地址（实际项目用 OSS/CDN）

package model

import "time"

// Product 商品模型
type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey"`                       // 主键
	Name        string    `json:"name" gorm:"size:100;not null"`              // 商品名称
	Description string    `json:"description" gorm:"size:500"`                // 商品描述
	Price       float64   `json:"price" gorm:"type:decimal(10,2);not null"`   // 价格，保留2位小数
	Stock       int       `json:"stock" gorm:"default:0"`                      // 库存
	ImageURL    string    `json:"image_url" gorm:"size:255"`                   // 商品图片地址
	CreatedAt   time.Time `json:"created_at"`                                  // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`                                  // 更新时间
}

// TableName 自定义表名
func (Product) TableName() string {
	return "products"
}
