package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"blog/model"
	"blog/pkg/cache"
	"blog/pkg/logger"
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
		Summary  string `json:"summary"`
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
		Summary: req.Summary,
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
func (h *ArticleHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	keyword := c.Query("keyword")

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
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	query.Count(&total)
	query.Order("created_at DESC").Offset(offset).Limit(size).Find(&articles)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  articles,
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

// ---------- 文章详情 ----------
func (h *ArticleHandler) Detail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "无效的ID"})
		return
	}

	cacheKey := fmt.Sprintf("article:detail:%d", id) //生成唯一的缓存键

	//---------- 1. 先尝试从 Redis 获取 ----------
	val, err := cache.Rdb.Get(cache.Ctx, cacheKey).Result()
	if err == nil {
		// 缓存命中！直接返回 JSON 数据
		logger.Log.Info("✅ 缓存命中", zap.String("key", cacheKey))
		var article model.Article
		if err := json.Unmarshal([]byte(val), &article); err == nil {
			// 异步原子递增阅读量（缓存里的 views 允许短暂滞后，最终一致）
			h.increaseViews(uint(id))
			c.JSON(http.StatusOK, gin.H{"code": 0, "data": article})
			return
		}
		// 如果 JSON 解析失败，就当缓存失败，继续查数据库
	}

	// ---------- 2. 缓存未命中，查询数据库 ----------
	logger.Log.Info("❌ 缓存未命中，查询数据库", zap.String("key", cacheKey))

	var article model.Article
	if err := h.DB.Preload("User").First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "文章不存在"})
		return
	}

	// 本次访问算一次阅读，写入缓存时带上 +1 后的值，避免缓存里的 views 一直滞后
	article.Views++

	// ---------- 3. 把查询结果写入 Redis（设置 5 分钟过期） ----------
	// 先把结构体转成 JSON 字符串
	jsonData, err := json.Marshal(article)
	if err != nil {
		logger.Log.Warn("序列化文章数据失败", zap.Error(err))
	} else {
		// 设置过期时间 5 分钟
		setErr := cache.Rdb.Set(cache.Ctx, cacheKey, jsonData, 5*time.Minute).Err()
		if setErr != nil {
			logger.Log.Warn("写入 Redis 缓存失败", zap.Error(setErr))
		} else {
			logger.Log.Info("✅ 写入 Redis 缓存成功", zap.String("key", cacheKey))
		}
	}

	// 增加阅读量（异步，不阻塞返回）
	// 用 SQL 原子递增（views = views + 1），并发下不会互相覆盖
	h.increaseViews(uint(id))

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": article})

}

// increaseViews 原子递增文章阅读量（异步调用）
func (h *ArticleHandler) increaseViews(id uint) {
	go func() {
		if err := h.DB.Model(&model.Article{}).Where("id = ?", id).
			UpdateColumn("views", gorm.Expr("views + 1")).Error; err != nil {
			logger.Log.Warn("递增阅读量失败", zap.Uint("id", id), zap.Error(err))
		}
	}()
}

// deleteCache 删除某篇文章的 Redis 缓存（更新/删除后调用，防止脏数据）
func (h *ArticleHandler) deleteCache(id int) {
	cacheKey := fmt.Sprintf("article:detail:%d", id)
	if err := cache.Rdb.Del(cache.Ctx, cacheKey).Err(); err != nil {
		logger.Log.Warn("删除 Redis 缓存失败", zap.String("key", cacheKey), zap.Error(err))
	} else {
		logger.Log.Info("✅ 已删除 Redis 缓存", zap.String("key", cacheKey))
	}
}

// ---------- 更新文章 ----------
func (h *ArticleHandler) Update(c *gin.Context) {
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
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Summary string `json:"summary"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "参数错误"})
		return
	}

	//更新字段
	updates := map[string]interface{}{
		"title":   req.Title,
		"content": req.Content,
		"summary": req.Summary,
	}

	if err := h.DB.Model(&article).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "更新失败：" + err.Error()})
		return
	}

	// 更新成功后删除 Redis 缓存，防止用户看到旧数据
	h.deleteCache(id)

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

	//权限检查： 只能删除自己的文章
	if article.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "无权删除此文章"})
		return
	}

	if err := h.DB.Delete(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "删除失败：" + err.Error()})
		return
	}

	// 删除成功后清理 Redis 缓存
	h.deleteCache(id)

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除成功"})

}
