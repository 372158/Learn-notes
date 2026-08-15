// internal/service/cart.go
// 购物车业务逻辑层
// 学习要点：购物车业务，Redis 与 MySQL 协作
//
// 业务流程：
// - 加购：查商品是否存在 -> 累加数量 -> 存 Redis
// - 查看购物车：读 Redis 商品列表 -> 批量查 MySQL 商品详情 -> 计算小计
// - 修改数量：直接覆盖 Redis 中的数量
// - 移除商品：从 Redis 删除字段

package service

import (
	"errors"

	"simple-ecommerce/internal/dao"
	"simple-ecommerce/internal/model"
)

// CartService 购物车业务服务
type CartService struct {
	cartDAO     *dao.CartDAO
	productDAO  *dao.ProductDAO
}

// NewCartService 创建 CartService 实例
// 学习要点：一个 service 可以依赖多个 DAO
func NewCartService() *CartService {
	return &CartService{
		cartDAO:    dao.NewCartDAO(),
		productDAO: dao.NewProductDAO(),
	}
}

// AddToCart 添加商品到购物车
// 学习要点：加购要累加数量，而非覆盖
func (s *CartService) AddToCart(userID, productID uint, quantity int) error {
	if quantity < 1 {
		return errors.New("数量必须大于 0")
	}

	// 1. 检查商品是否存在
	product, err := s.productDAO.GetByID(productID)
	if err != nil {
		return errors.New("商品不存在")
	}

	// 2. 检查库存
	// 学习要点：加购时校验库存，避免加购数量超过库存
	existQty, err := s.cartDAO.GetQuantity(userID, productID)
	if err != nil {
		existQty = 0 // 不存在表示还没加过，数量为0
	}
	if existQty+quantity > product.Stock {
		return errors.New("购物车数量超过库存")
	}

	// 3. 累加数量并保存
	// 学习要点：先读原有数量，加上新数量，再写回
	return s.cartDAO.UpdateQuantity(userID, productID, existQty+quantity)
}

// GetCart 查看购物车
// 学习要点：Redis + MySQL 协作
// 1. 从 Redis 读购物车（只有 product_id 和 quantity）
// 2. 根据 product_id 批量从 MySQL 查商品详情
// 3. 组装成 CartItem 返回前端
func (s *CartService) GetCart(userID uint) ([]model.CartItem, float64, error) {
	// 1. 从 Redis 读购物车
	items, err := s.cartDAO.GetAllItems(userID)
	if err != nil {
		return nil, 0, err
	}
	if len(items) == 0 {
		return []model.CartItem{}, 0, nil
	}

	// 2. 查商品详情并组装
	var cartItems []model.CartItem
	var totalPrice float64
	for productID, quantity := range items {
		product, err := s.productDAO.GetByID(productID)
		if err != nil {
			// 商品可能已下架，跳过
			continue
		}
		subtotal := product.Price * float64(quantity)
		cartItems = append(cartItems, model.CartItem{
			ProductID:  productID,
			Name:       product.Name,
			Price:      product.Price,
			Quantity:   quantity,
			ImageURL:   product.ImageURL,
			TotalPrice: subtotal,
		})
		totalPrice += subtotal
	}

	return cartItems, totalPrice, nil
}

// UpdateQuantity 修改购物车商品数量
func (s *CartService) UpdateQuantity(userID, productID uint, quantity int) error {
	if quantity < 1 {
		return errors.New("数量必须大于 0")
	}

	// 检查商品是否存在及库存
	product, err := s.productDAO.GetByID(productID)
	if err != nil {
		return errors.New("商品不存在")
	}
	if quantity > product.Stock {
		return errors.New("数量超过库存")
	}

	return s.cartDAO.UpdateQuantity(userID, productID, quantity)
}

// RemoveItem 移除购物车商品
func (s *CartService) RemoveItem(userID, productID uint) error {
	return s.cartDAO.RemoveItem(userID, productID)
}

// GetCartItemsForOrder 获取购物车内容用于下单
// 学习要点：内部方法，供 order service 调用，不返回给前端
func (s *CartService) GetCartItemsForOrder(userID uint) (map[uint]int, error) {
	return s.cartDAO.GetAllItems(userID)
}

// ClearCart 清空购物车
// 学习要点：下单成功后清空购物车
func (s *CartService) ClearCart(userID uint) error {
	return s.cartDAO.ClearCart(userID)
}
