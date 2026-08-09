package handler

import (
	"net/http"
	"strconv"

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


// ---------- 文章列表（分页 + 搜索） ----------
func (h *ArticleHandler) List(c *gin.Context)  {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	keyword := c.Query("kewword")

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 10
	}
	offset := (page - 1) * size
	var articles []model.Article
	var total int64

	query := h.DB.Model(&model.Article{}).Preload("User")

	//标题模糊查询
	if keyword != "" {
		query = query.Where("title LIKE ?", "%" + keyword + "%")
	}

	query.Count(&total)
	query.Order("created_at DESC").Offset(offset).Limit(size).Find(&articles)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list": articles,
			"total": total,
			"page": page,
			"size": size,
		},
	})
}


//---------- 文章详情 ----------
func (h *ArticleHandler) Detail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "无效的ID"})
		return
	}

	var article model.Article
	if err := h.DB.Preload("User").First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "文章不存在"})
		return
	}

	//增加阅读量（异步， 不阻塞返回）
	go func() {
		h.DB.Model(&article).UpdateColumn("views", article.Views + 1)
	}()

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": article})

}


// ---------- 更新文章 ----------
func (h *ArticleHandler) Updata(c *gin.Context) {
	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未登录"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "无效的ID"})
		return
	}

	//查询文章
	var article model.Article
	if err := h.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "文章不存在"})
		return
	}

	//权限检查：只能修改自己的文章
	if article.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "无权修改此文章"})
		return
	}

	//绑定请求参数
	var req struct {
		Title		string	`json:"title" binding:"required"`
		Content		string	`json:"content" binding:"required"`
		Summary		string	`json:"summary"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "参数错误"})
		return
	}

	//更新字段
	updates := map[string]interface{} {
		"title": req.Title,
		"content": req.Content,
		"summary":  req.Summary,
	}

	h.DB.Model(&article).Updates(updates)

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "更新成功", "data": article})
	
}


// 删除文章
	func (h *ArticleHandler) Delete(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未登录"})
			return
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "无效的ID"})
			return
		}

		var article model.Article
		if err := h.DB.First(&article, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "文章不存在"})
			return
		}

		//权限检查： 只能三处自己的文章
		if article.UserID != userID.(uint) {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "无权删除此文章"})
			return
		}

		h.DB.Delete(&article)

		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除成功"})
		
	}
