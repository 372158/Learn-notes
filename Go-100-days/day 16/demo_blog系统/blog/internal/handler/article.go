package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"blog/model"
)

type ArticleHandler struct {
	DB *gorm.DB
}

func NewArticleHandler(db *gorm.DB) *ArticleHandler {
	return &ArticleHandler{DB: db}
}

// CreateArticle 创建文章（需要登录）
func (h *ArticleHandler) Create(c *gin.Context) {
	// 1. 从上下文获取当前登录用户ID(由 JWT 中间件注入)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未登录",
		})
		return
	}

	// 2. 绑定请求参数
	var req struct {
		Title    string `json:"title" binding:"required"`
		Content  string `json:"content" binding:"required"`
		Sunmmary string `json:"summary"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误:" + err.Error(),
		})
		return
	}

	// 3. 创建文章
	article := model.Article{
		Title:   req.Title,
		Content: req.Content,
		Summary: req.Sunmmary,
		UserID:  userID.(uint), //类型断言
	}

	if err := h.DB.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建文章失败：" + err.Error(),
		})
		return
	}

	// 4. 返回结果
	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"msg":  "创建成功",
		"data": article,
	})
}
