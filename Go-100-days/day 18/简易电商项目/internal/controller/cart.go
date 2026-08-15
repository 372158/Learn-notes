// internal/controller/cart.go
// 购物车控制器层
// 学习要点：购物车接口需要登录，从 context 获取 user_id
//
// 接口设计：
// - POST   /api/v1/cart            加入购物车
// - GET    /api/v1/cart            查看购物车
// - PUT    /api/v1/cart/:product_id 修改数量
// - DELETE /api/v1/cart/:product_id 移除商品

package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"simple-ecommerce/internal/service"
	"simple-ecommerce/pkg/response"
)

// CartController 购物车控制器
type CartController struct {
	cartService *service.CartService
}

// NewCartController 创建 CartController 实例
func NewCartController() *CartController {
	return &CartController{
		cartService: service.NewCartService(),
	}
}

// AddToCartRequest 加购请求参数
type AddToCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"` // 商品 ID
	Quantity  int  `json:"quantity" binding:"required,min=1"` // 数量（最少1）
}

// AddToCart 加入购物车
// 接口：POST /api/v1/cart
func (ctrl *CartController) AddToCart(c *gin.Context) {
	// 学习要点：从 context 获取当前登录用户 ID（JWT 中间件设置）
	userID := c.MustGet("user_id").(uint)

	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.cartService.AddToCart(userID, req.ProductID, req.Quantity); err != nil {
		response.Error(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "已加入购物车", nil)
}

// GetCart 查看购物车
// 接口：GET /api/v1/cart
func (ctrl *CartController) GetCart(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	items, totalPrice, err := ctrl.cartService.GetCart(userID)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"items":       items,
		"total_price": totalPrice,
		"count":       len(items),
	})
}

// UpdateQuantityRequest 修改数量请求参数
type UpdateQuantityRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

// UpdateQuantity 修改购物车商品数量
// 接口：PUT /api/v1/cart/:product_id
func (ctrl *CartController) UpdateQuantity(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		response.Error(c, "商品 ID 格式错误")
		return
	}

	var req UpdateQuantityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.cartService.UpdateQuantity(userID, uint(productID), req.Quantity); err != nil {
		response.Error(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "修改成功", nil)
}

// RemoveItem 移除购物车商品
// 接口：DELETE /api/v1/cart/:product_id
func (ctrl *CartController) RemoveItem(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		response.Error(c, "商品 ID 格式错误")
		return
	}

	if err := ctrl.cartService.RemoveItem(userID, uint(productID)); err != nil {
		response.Error(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "已移除", nil)
}
