// internal/service/order.go
// 订单业务逻辑层
// 学习要点：下单核心业务流程编排
//
// 下单完整流程：
// 1. 读取购物车内容
// 2. 校验购物车非空
// 3. 查商品详情，生成订单项（锁定价格快照）
// 4. 计算订单总金额
// 5. 生成订单号
// 6. 事务创建订单 + 订单项 + 扣库存
// 7. 清空购物车
// 8. 返回订单信息

package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"simple-ecommerce/internal/dao"
	"simple-ecommerce/internal/model"
)

// OrderService 订单业务服务
type OrderService struct {
	orderDAO   *dao.OrderDAO
	cartDAO    *dao.CartDAO
	productDAO *dao.ProductDAO
	cartService *CartService
}

// NewOrderService 创建 OrderService 实例
// 学习要点：订单服务依赖多个 DAO 和购物车服务
func NewOrderService() *OrderService {
	return &OrderService{
		orderDAO:    dao.NewOrderDAO(),
		cartDAO:     dao.NewCartDAO(),
		productDAO:  dao.NewProductDAO(),
		cartService: NewCartService(),
	}
}

// generateOrderNo 生成订单号
// 学习要点：订单号生成规则
// 格式：年月日时分秒 + 6位随机数，例如 20260814153012 123456
// 这样保证：可读性（能看出下单时间）+ 唯一性（随机数）
func generateOrderNo() string {
	now := time.Now()
	// 时间部分：20260814153012
	timeStr := now.Format("20060102150405")
	// 随机部分：6位随机数
	randNum := rand.Intn(1000000)
	return fmt.Sprintf("%s%06d", timeStr, randNum)
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(userID uint) (*model.Order, error) {
	// 1. 读取购物车内容
	cartItems, err := s.cartDAO.GetAllItems(userID)
	if err != nil {
		return nil, errors.New("读取购物车失败: " + err.Error())
	}
	if len(cartItems) == 0 {
		return nil, errors.New("购物车为空")
	}

	// 2. 生成订单项 + 计算总金额
	var orderItems []model.OrderItem
	var totalPrice float64
	for productID, quantity := range cartItems {
		product, err := s.productDAO.GetByID(productID)
		if err != nil {
			return nil, fmt.Errorf("商品(ID:%d)不存在", productID)
		}

		// 学习要点：订单项存储商品快照（名称、价格）
		// 即使商品后续改名或改价，历史订单不受影响
		subtotal := product.Price * float64(quantity)
		orderItems = append(orderItems, model.OrderItem{
			ProductID:   productID,
			ProductName: product.Name,
			Price:       product.Price,
			Quantity:    quantity,
			TotalPrice:  subtotal,
		})
		totalPrice += subtotal
	}

	// 3. 组装订单
	order := &model.Order{
		OrderNo:    generateOrderNo(),
		UserID:     userID,
		TotalPrice: totalPrice,
		Status:     model.OrderStatusPending, // 待支付
		Items:      orderItems,
	}

	// 4. 事务创建订单（含扣库存）
	// 学习要点：传入 productDAO 用于事务内扣库存
	if err := s.orderDAO.CreateOrder(order, s.productDAO); err != nil {
		if errors.Is(err, dao.ErrInsufficientStock) {
			return nil, errors.New("库存不足")
		}
		return nil, errors.New("创建订单失败: " + err.Error())
	}

	// 5. 清空购物车
	// 学习要点：下单成功后清空购物车
	// 注意：清购物车不在事务内，因为 Redis 不支持 MySQL 事务
	// 即使清购物车失败，订单已创建成功，不影响订单
	if err := s.cartDAO.ClearCart(userID); err != nil {
		// 仅记录日志，不影响订单
		fmt.Printf("警告：清空购物车失败: %v\n", err)
	}

	return order, nil
}

// GetOrder 查询订单详情
func (s *OrderService) GetOrder(id uint, userID uint) (*model.Order, error) {
	order, err := s.orderDAO.GetByID(id, userID)
	if err != nil {
		return nil, errors.New("订单不存在")
	}
	return order, nil
}

// OrderListResult 订单列表结果
type OrderListResult struct {
	List  []model.Order `json:"list"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Size  int           `json:"size"`
}

// GetOrderList 查询用户订单列表
func (s *OrderService) GetOrderList(userID uint, page, size int) (*OrderListResult, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	orders, total, err := s.orderDAO.GetListByUserID(userID, page, size)
	if err != nil {
		return nil, err
	}

	return &OrderListResult{
		List:  orders,
		Total: total,
		Page:  page,
		Size:  size,
	}, nil
}

// PayOrder 支付订单（模拟）
// 学习要点：模拟支付，实际项目对接支付宝/微信支付
func (s *OrderService) PayOrder(id uint, userID uint) error {
	order, err := s.orderDAO.GetByID(id, userID)
	if err != nil {
		return errors.New("订单不存在")
	}

	if order.Status != model.OrderStatusPending {
		return errors.New("订单状态不允许支付")
	}

	return s.orderDAO.UpdateStatus(id, model.OrderStatusPaid)
}
