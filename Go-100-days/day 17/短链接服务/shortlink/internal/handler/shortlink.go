package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"shortlink/internal/service"
)

type ShortLinkHandler struct {
	Service *service.ShortLinkService
}

func NewShortLinkHandler(svc *service.ShortLinkService) *ShortLinkHandler {
	return &ShortLinkHandler{Service: svc}
}

// Create 生成短链接
func (h *ShortLinkHandler) Create(c *gin.Context) {
	var req struct {
		URL       string `json:"url" binding:"required"`
		ExpiredAt string `json:"expired_at"` // 可选，格式 "2026-12-31T23:59:59Z"
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请提供 URL"})
		return
	}

	shortLink, err := h.Service.CreateShortLink(req.URL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"short_code": shortLink.ShortCode,
			"short_url":  "http://localhost:8080/" + shortLink.ShortCode,
			"original":   shortLink.OriginalURL,
		},
	})
}

// Redirect 重定向到原始 URL
func (h *ShortLinkHandler) Redirect(c *gin.Context) {
	shortCode := c.Param("code")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "缺少短码"})
		return
	}

	originalURL, err := h.Service.GetOriginalURL(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": err.Error()})
		return
	}

	c.Redirect(http.StatusMovedPermanently, originalURL)
}