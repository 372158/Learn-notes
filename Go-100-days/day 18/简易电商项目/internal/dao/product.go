// internal/dao/product.go
// 商品数据访问层
// 学习要点：商品 CRUD + 分页查询
//
// 分页查询核心：
// - Offset: 跳过多少条（偏移量）
// - Limit: 取多少条
// 公式：offset = (page - 1) * size
// 例如第2页每页10条：offset = (2-1)*10 = 10，跳过前10条取后10条

package dao

import (
	"simple-ecommerce/internal/model"
	"simple-ecommerce/initialize"

	"gorm.io/gorm"
)

// ProductDAO 商品数据访问对象
type ProductDAO struct{}

// NewProductDAO 创建 ProductDAO 实例
func NewProductDAO() *ProductDAO {
	return &ProductDAO{}
}

// Create 创建商品
// 学习要点：db.Create 插入数据，回填自增 ID
func (dao *ProductDAO) Create(product *model.Product) error {
	return initialize.DB.Create(product).Error
}

// GetByID 根据 ID 查询商品
// 学习要点：db.First 查询单条，第二个参数是主键
func (dao *ProductDAO) GetByID(id uint) (*model.Product, error) {
	var product model.Product
	err := initialize.DB.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Update 更新商品
// 学习要点：db.Model().Updates() 更新非零字段
// - Model 指定操作的模型和主键
// - Updates 只更新结构体中非零值的字段
func (dao *ProductDAO) Update(product *model.Product) error {
	return initialize.DB.Model(product).Updates(product).Error
}

// Delete 删除商品
// 学习要点：db.Delete 软删除（如果模型有 DeletedAt 字段）或硬删除
// 这里 Product 没有 DeletedAt，是硬删除
func (dao *ProductDAO) Delete(id uint) error {
	return initialize.DB.Delete(&model.Product{}, id).Error
}

// GetList 分页查询商品列表
// 学习要点：分页查询三件套
// 1. Count: 查总数（用于前端显示总页数）
// 2. Offset + Limit: 分页取数据
// 3. Order: 排序（按创建时间倒序，新的在前）
//
// 参数：
// - page: 页码（从1开始）
// - size: 每页条数
// 返回：
// - 商品列表、总数、错误
func (dao *ProductDAO) GetList(page, size int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	// 先查总数
	// 学习要点：Model 指定表，Count 统计数量
	initialize.DB.Model(&model.Product{}).Count(&total)

	// 再分页查询
	// 学习要点：链式调用，按顺序执行
	// - Offset((page-1)*size): 跳过前面页的数据
	// - Limit(size): 取 size 条
	// - Order("created_at DESC"): 按创建时间倒序
	offset := (page - 1) * size
	err := initialize.DB.Offset(offset).Limit(size).Order("created_at DESC").Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// UpdateStock 扣减库存
// 学习要点：用 SQL 表达式扣减库存，避免并发问题
// db.Model(&Product{}).Where("id = ? AND stock >= ?", id, num).Update("stock", gorm.Expr("stock - ?", num))
// 这样能保证库存不足时更新失败（返回影响行数为0）
func (dao *ProductDAO) UpdateStock(id uint, quantity int) error {
	result := initialize.DB.Model(&model.Product{}).
		Where("id = ? AND stock >= ?", id, quantity).
		Update("stock", gorm.Expr("stock - ?", quantity))
	if result.Error != nil {
		return result.Error
	}
	// 学习要点：RowsAffected 返回受影响的行数
	// 如果库存不足，没有记录被更新，RowsAffected 为 0
	if result.RowsAffected == 0 {
		return ErrInsufficientStock
	}
	return nil
}
