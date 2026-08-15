// internal/router/router.go
// 路由注册模块
// 学习要点：Gin 路由分组，区分公开接口和鉴权接口
//
// 路由分组设计：
// 1. 公开组（无需登录）：注册、登录、商品列表/详情
// 2. 鉴权组（需登录）：商品增删改、购物车、订单
//
// 为什么商品列表放公开组：
// 浏览商品不需要登录，降低使用门槛（实际电商都这样设计）

package router

import (
	"github.com/gin-gonic/gin"

	"simple-ecommerce/internal/controller"
	"simple-ecommerce/internal/middleware"
)

// SetupRouter 初始化路由
// 学习要点：集中注册所有路由，返回配置好的 gin.Engine
func SetupRouter() *gin.Engine {
	// 学习要点：gin.Default() 创建带 Logger 和 Recovery 中间件的路由引擎
	// - Logger: 打印请求日志
	// - Recovery: panic 自动恢复，返回 500
	r := gin.Default()

	// 创建各 controller 实例
	userCtrl := controller.NewUserController()
	productCtrl := controller.NewProductController()
	cartCtrl := controller.NewCartController()
	orderCtrl := controller.NewOrderController()

	// API 版本分组（便于未来升级 v2）
	// 学习要点：Group 创建路由分组，路径前缀统一
	api := r.Group("/api/v1")

	// ============ 公开接口（无需登录）============
	api.POST("/register", userCtrl.Register) // 注册
	api.POST("/login", userCtrl.Login)       // 登录

	// 商品浏览（公开）
	api.GET("/products", productCtrl.GetList)     // 商品列表
	api.GET("/products/:id", productCtrl.GetByID) // 商品详情

	// ============ 鉴权接口（需登录）============
	// 学习要点：Use 注册中间件，该组下所有路由都会先经过 JWT 鉴权
	auth := api.Group("/")
	auth.Use(middleware.JWTAuth())
	{
		// 商品管理（需登录）
		auth.POST("/products", productCtrl.Create)       // 创建商品
		auth.PUT("/products/:id", productCtrl.Update)    // 更新商品
		auth.DELETE("/products/:id", productCtrl.Delete) // 删除商品

		// 购物车
		auth.POST("/cart", cartCtrl.AddToCart)                  // 加入购物车
		auth.GET("/cart", cartCtrl.GetCart)                     // 查看购物车
		auth.PUT("/cart/:product_id", cartCtrl.UpdateQuantity)  // 修改数量
		auth.DELETE("/cart/:product_id", cartCtrl.RemoveItem)   // 移除商品

		// 订单
		auth.POST("/orders", orderCtrl.CreateOrder)        // 创建订单
		auth.GET("/orders", orderCtrl.GetOrderList)        // 订单列表
		auth.GET("/orders/:id", orderCtrl.GetOrder)        // 订单详情
		auth.PUT("/orders/:id/pay", orderCtrl.PayOrder)    // 支付订单
	}

	return r
}
