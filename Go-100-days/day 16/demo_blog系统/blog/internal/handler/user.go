package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"blog/internal/middleware"
	"blog/model"
)

// ---------- UserHandler ----------
type UserHandler struct {
	DB *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

// ---------- 注册 ----------

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest

	// 1. 绑定 JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 2. 检查用户名是否已被占用
	var existUser model.User
	if err := h.DB.Where("username = ?", req.Username).First(&existUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"code": 409,
			"msg":  "用户名已被注册",
		})
		return
	}

	// 3. 检查邮箱是否已被占用
	if err := h.DB.Where("email = ?", req.Email).First(&existUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"code": 409,
			"msg":  "邮箱已被注册",
		})
		return
	}

	// 4. 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "加密密码失败",
		})
		return
	}

	// 5. 创建用户
	user := model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Role:     "user",
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建用户失败: " + err.Error(),
		})
		return
	}

	// 6. 返回结果（不返回密码）
	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"msg":  "注册成功",
		"data": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"role":       user.Role,
			"created_at": user.CreatedAt,
		},
	})
}

// ---------- 登录 ----------

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest

	// 1. 绑定 JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 2. 查询用户
	var user model.User
	if err := h.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		// 用户名不存在
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "用户名或密码错误",
		})
		return
	}

	// 3. 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "用户名或密码错误",
		})
		return
	}

	// 4. 生成 JWT Token
	secret := viper.GetString("jwt.secret")
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "JWT 密钥未配置",
		})
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "生成 Token 失败: " + err.Error(),
		})
		return
	}

	// 5. 返回 Token
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "登录成功",
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"role":     user.Role,
			},
		},
	})
}
