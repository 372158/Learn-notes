// internal/dao/order.go
// 订单数据访问层
// 学习要点：事务（Transaction）的使用
//
// 为什么创建订单要用事务：
// 下单涉及多步操作，任一步失败都要全部回滚：
// 1. 创建订单主表记录
// 2. 创建多个订单项记录
// 3. 扣减商品库存
// 4. 清空购物车
// 如果中途失败（如库存不足），前面已执行的必须撤销，否则数据不一致
//
// GORM 事务用法：
// db.Transaction(func(tx *gorm.DB) error {
//     // 在事务中执行操作，用 tx 而非 db
//     // 返回 nil 提交事务，返回 error 回滚事务
// })

package dao

import (
	"simple-ecommerce/internal/model"
	"simple-ecommerce/initialize"

	"gorm.io/gorm"
)

// OrderDAO 订单数据访问对象
type OrderDAO struct{}

// NewOrderDAO 创建 OrderDAO 实例
func NewOrderDAO() *OrderDAO {
	return &OrderDAO{}
}

// CreateOrder 创建订单（事务）
// 学习要点：事务保证原子性，要么全部成功，要么全部失败
func (dao *OrderDAO) CreateOrder(order *model.Order, productDAO *ProductDAO) error {
	// db.Transaction 自动管理事务
	// - 闭包返回 nil -> 自动 Commit（提交）
	// - 闭包返回 error -> 自动 Rollback（回滚）
	return initialize.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 创建订单主表
		// 学习要点：事务中用 tx 而非 db
		if err := tx.Create(order).Error; err != nil {
			return err // 返回 error 会自动回滚
		}

		// 2. 创建订单项 + 扣减库存
		for _, item := range order.Items {
			item.OrderID = order.ID // 关联订单 ID
			if err := tx.Create(&item).Error; err != nil {
				return err
			}

			// 扣减库存（在事务中执行，保证库存扣减和订单创建原子性）
			// 学习要点：用 tx 操作，确保和订单创建在同一事务
			if err := dao.updateStockInTx(tx, item.ProductID, item.Quantity); err != nil {
				return err // 库存不足会回滚整个事务
			}
		}

		return nil // 返回 nil 提交事务
	})
}

// updateStockInTx 在事务中扣减库存
// 学习要点：事务内库存扣减，用 tx 保证和订单操作同事务
func (dao *OrderDAO) updateStockInTx(tx *gorm.DB, productID uint, quantity int) error {
	result := tx.Model(&model.Product{}).
		Where("id = ? AND stock >= ?", productID, quantity).
		Update("stock", gorm.Expr("stock - ?", quantity))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrInsufficientStock
	}
	return nil
}

// GetByID 查询订单详情（含订单项）
// 学习要点：Preload 预加载关联数据
// - Preload("Items"): 自动查询关联的订单项
func (dao *OrderDAO) GetByID(id uint, userID uint) (*model.Order, error) {
	var order model.Order
	err := initialize.DB.
		Where("id = ? AND user_id = ?", id, userID).
		Preload("Items").
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetListByUserID 分页查询用户订单
// 学习要点：关联查询分页，先查主表分页再 Preload
func (dao *OrderDAO) GetListByUserID(userID uint, page, size int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	initialize.DB.Model(&model.Order{}).Where("user_id = ?", userID).Count(&total)

	offset := (page - 1) * size
	err := initialize.DB.
		Where("user_id = ?", userID).
		Preload("Items").
		Offset(offset).Limit(size).
		Order("created_at DESC").
		Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// UpdateStatus 更新订单状态
func (dao *OrderDAO) UpdateStatus(id uint, status int) error {
	return initialize.DB.Model(&model.Order{}).Where("id = ?", id).Update("status", status).Error
}
