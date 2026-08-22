package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"llm-chat-db/auth"
)

// Auth 校验请求头中的 `Authorization: Bearer <token>`：
// 校验通过把用户名放进上下文并放行，否则 401 并中断。
// 用法：在受保护路由前挂上，如 r.GET("/chat", middleware.Auth(secret), handler)
func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") { // 请求头格式要求：Bearer + 空格 + token
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少 token 或格式错误"})
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(h, "Bearer ")
		claims, err := auth.Parse(secret, tokenStr)
		if err != nil {
			// 签名被篡改、已过期、格式错都走到这
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效或过期 token"})
			c.Abort()
			return
		}
		c.Set("user", claims.User) // 后面 handler/service 可用 c.Get("user")
		c.Next()
	}
}
