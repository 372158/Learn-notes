// internal/controller/order.go
// 订单控制器层
// 学习要点：订单接口均需登录，从 context 获取 user_id
//
// 接口设计：
// - POST   /api/v1/orders      创建订单（从购物车结算）
// - GET    /api/v1/orders      我的订单列表
// - GET    /api/v1/orders/:id  订单详情
// - PUT    /api/v1/orders/:id/pay  支付订单（模拟）

package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"simple-ecommerce/internal/service"
	"simple-ecommerce/pkg/response"
)

// OrderController 订单控制器
type OrderController struct {
	orderService *service.OrderService
}

// NewOrderController 创建 OrderController 实例
func NewOrderController() *OrderController {
	return &OrderController{
		orderService: service.NewOrderService(),
	}
}

// CreateOrder 创建订单
// 接口：POST /api/v1/orders
// 学习要点：从购物车结算下单，不需要传参数
func (ctrl *OrderController) CreateOrder(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	order, err := ctrl.orderService.CreateOrder(userID)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "下单成功", order)
}

// GetOrder 查询订单详情
// 接口：GET /api/v1/orders/:id
func (ctrl *OrderController) GetOrder(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, "订单 ID 格式错误")
		return
	}

	order, err := ctrl.orderService.GetOrder(uint(id), userID)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, order)
}

// GetOrderList 查询我的订单列表
// 接口：GET /api/v1/orders?page=1&size=10
func (ctrl *OrderController) GetOrderList(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	result, err := ctrl.orderService.GetOrderList(userID, page, size)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, result)
}

// PayOrder 支付订单（模拟）
// 接口：PUT /api/v1/orders/:id/pay
func (ctrl *OrderController) PayOrder(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, "订单 ID 格式错误")
		return
	}

	if err := ctrl.orderService.PayOrder(uint(id), userID); err != nil {
		response.Error(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "支付成功", nil)
}
