package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"agent-demo/auth"
)

// Auth 校验请求头 Authorization: Bearer <token>，通过后把用户名写进上下文
func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "缺少 token"})
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := auth.Parse(secret, tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token 无效或已过期"})
			return
		}
		c.Set("user", claims.User)
		c.Next()
	}
}