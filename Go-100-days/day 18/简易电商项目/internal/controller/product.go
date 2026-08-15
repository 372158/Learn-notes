// internal/controller/product.go
// 商品控制器层
// 学习要点：RESTful 风格的 CRUD 接口设计
//
// RESTful 接口设计规范：
// - GET    /products      获取列表
// - GET    /products/:id  获取详情
// - POST   /products      创建
// - PUT    /products/:id  更新
// - DELETE /products/:id  删除
//
// HTTP 方法语义：
// - GET: 查询（安全、幂等）
// - POST: 创建
// - PUT: 更新（幂等）
// - DELETE: 删除（幂等）

package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"simple-ecommerce/internal/model"
	"simple-ecommerce/internal/service"
	"simple-ecommerce/pkg/response"
)

// ProductController 商品控制器
type ProductController struct {
	productService *service.ProductService
}

// NewProductController 创建 ProductController 实例
func NewProductController() *ProductController {
	return &ProductController{
		productService: service.NewProductService(),
	}
}

// CreateProductRequest 创建商品请求参数
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=100"`   // 商品名称
	Description string  `json:"description" binding:"max=500"`           // 商品描述
	Price       float64 `json:"price" binding:"required,gt=0"`           // 价格（必须大于0）
	Stock       int     `json:"stock" binding:"gte=0"`                    // 库存（>=0）
	ImageURL    string  `json:"image_url" binding:"max=255"`             // 图片地址
}

// Create 创建商品
// 接口：POST /api/v1/products
func (ctrl *ProductController) Create(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "参数错误: "+err.Error())
		return
	}

	product := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
	}

	if err := ctrl.productService.CreateProduct(product); err != nil {
		response.Error(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "创建成功", product)
}

// GetByID 查询商品详情
// 接口：GET /api/v1/products/:id
func (ctrl *ProductController) GetByID(c *gin.Context) {
	// 学习要点：c.Param 获取路径参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, "商品 ID 格式错误")
		return
	}

	product, err := ctrl.productService.GetProduct(uint(id))
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, product)
}

// GetList 分页查询商品列表
// 接口：GET /api/v1/products?page=1&size=10
func (ctrl *ProductController) GetList(c *gin.Context) {
	// 学习要点：c.Query 获取查询参数（返回字符串，需转换）
	// c.DefaultQuery 带默认值
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	result, err := ctrl.productService.GetProductList(page, size)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, result)
}

// UpdateProductRequest 更新商品请求参数
type UpdateProductRequest struct {
	Name        string  `json:"name" binding:"omitempty,min=1,max=100"`
	Description string  `json:"description" binding:"omitempty,max=500"`
	Price       float64 `json:"price" binding:"omitempty,gt=0"`
	Stock       int     `json:"stock" binding:"omitempty,gte=0"`
	ImageURL    string  `json:"image_url" binding:"omitempty,max=255"`
}

// Update 更新商品
// 接口：PUT /api/v1/products/:id
func (ctrl *ProductController) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, "商品 ID 格式错误")
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "参数错误: "+err.Error())
		return
	}

	// 学习要点：更新时用指针字段区分"不更新"和"更新为零值"
	// 这里简化处理，用结构体直接更新
	product := &model.Product{
		ID:          uint(id),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
	}

	if err := ctrl.productService.UpdateProduct(product); err != nil {
		response.Error(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "更新成功", nil)
}

// Delete 删除商品
// 接口：DELETE /api/v1/products/:id
func (ctrl *ProductController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, "商品 ID 格式错误")
		return
	}

	if err := ctrl.productService.DeleteProduct(uint(id)); err != nil {
		response.Error(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "删除成功", nil)
}
