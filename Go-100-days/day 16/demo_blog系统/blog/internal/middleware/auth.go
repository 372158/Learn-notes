package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Auth 认证中间件，验证 JWT Token

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. 从 Header 获取 Token
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg": "未提供 Token",
			})
			ctx.Abort()
			return 
		}

		// 2. 解析 Token （格式： Bearer <token>）
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":	"Token 格式错误，请使用 Bearer <token>",
			})
			ctx.Abort()
			return
		}

		// 3. 验证 Token
		claims , err := ValidateToken(parts[1])
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg": "Token 无效或已过期" + err.Error(),
			})
			ctx.Abort()
			return 
		}

		// 4. 将用户信息存入上下文， 供后续 Handler 使用
		ctx.Set("user_id", claims.UserID)
		ctx.Set("username", claims.Username)
		ctx.Set("role", claims.Role)

		ctx.Next()
	}
}