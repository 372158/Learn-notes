package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"llm-chat-db/auth"
)

// LoginHandler 处理登录并签发 JWT
type LoginHandler struct {
	Secret string
	TTL    time.Duration
}

func NewLoginHandler(secret string, ttl time.Duration) *LoginHandler {
	return &LoginHandler{Secret: secret, TTL: ttl}
}

// Login 校验用户名密码，通过则返回 token。演示用写死账号，
// 真实项目这里应查数据库/其他认证源做校验。
func (h *LoginHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if req.Username != "admin" || req.Password != "123456" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}
	token, err := auth.Generate(h.Secret, req.Username, h.TTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "签发 token 失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "expires_in": int(h.TTL.Seconds())})
}
