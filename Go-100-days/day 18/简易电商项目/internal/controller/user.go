// internal/controller/user.go
// 用户控制器层
// 学习要点：Controller 层负责接收 HTTP 请求、参数校验、调用 service、返回响应
//
// 分层职责：
// - controller: HTTP 相关（参数绑定、响应返回）
// - service: 业务逻辑（密码加密、校验）
// - dao: 数据库操作（增删改查）
//
// 请求处理流程：
// 1. c.ShouldBindJSON 绑定 JSON 参数到结构体
// 2. binding tag 自动校验（required、min、max 等）
// 3. 调用 service 处理业务
// 4. 用 response 包统一返回结果

package controller

import (
	"github.com/gin-gonic/gin"

	"simple-ecommerce/internal/service"
	"simple-ecommerce/pkg/response"
)

// UserController 用户控制器
type UserController struct {
	userService *service.UserService
}

// NewUserController 创建 UserController 实例
func NewUserController() *UserController {
	return &UserController{
		userService: service.NewUserService(),
	}
}

// RegisterRequest 注册请求参数
// 学习要点：用结构体定义请求参数，加 binding tag 自动校验
// - required: 必填
// - min: 最小长度
// - max: 最大长度
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`  // 用户名 3-50 字符
	Password string `json:"password" binding:"required,min=6,max=20"`  // 密码 6-20 字符
	Phone    string `json:"phone" binding:"omitempty,len=11"`          // 手机号 11 位，可选
}

// Register 用户注册接口
// 接口：POST /api/v1/register
func (ctrl *UserController) Register(c *gin.Context) {
	// 1. 绑定并校验参数
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "参数错误: "+err.Error())
		return
	}

	// 2. 调用 service 处理业务
	user, err := ctrl.userService.Register(req.Username, req.Password, req.Phone)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	// 3. 返回结果
	// 学习要点：不返回密码字段（model 中 json:"-" 已处理）
	response.SuccessWithMsg(c, "注册成功", gin.H{
		"id":       user.ID,
		"username": user.Username,
		"phone":    user.Phone,
	})
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 用户登录接口
// 接口：POST /api/v1/login
func (ctrl *UserController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "参数错误: "+err.Error())
		return
	}

	token, err := ctrl.userService.Login(req.Username, req.Password)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	// 学习要点：返回 token，前端保存后每次请求放在 Authorization 头
	response.SuccessWithMsg(c, "登录成功", gin.H{
		"token": token,
	})
}
