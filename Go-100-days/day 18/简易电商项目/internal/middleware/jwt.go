// internal/middleware/jwt.go
// JWT 鉴权中间件
// 学习要点：中间件是 Gin 的核心概念，能在请求到达 controller 前后做统一处理
//
// 为什么用中间件做鉴权：
// 1. 统一处理：所有需要登录的接口共用一套鉴权逻辑
// 2. 解耦：controller 不用重复写鉴权代码
// 3. 灵活：通过路由分组决定哪些接口需要鉴权
//
// JWT 鉴权流程：
// 1. 从请求头 Authorization 提取 token
// 2. 解析验证 token
// 3. 验证通过：把用户信息存入 gin.Context，供后续使用
// 4. 验证失败：返回 401 错误

package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"simple-ecommerce/pkg/jwt"
	"simple-ecommerce/pkg/response"
)

// JWTAuth JWT 鉴权中间件
// 学习要点：中间件函数签名 func(*gin.Context)，返回 gin.HandlerFunc
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 学习要点：从请求头获取 token
		// 约定格式：Authorization: Bearer <token>
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ErrorWithCode(c, response.CodeUnauthorized, "请求未携带 token")
			c.Abort() // 学习要点：Abort 阻止后续中间件和 controller 执行
			return
		}

		// 学习要点：按空格分割，取第二部分（token）
		// "Bearer abc123" -> ["Bearer", "abc123"]
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.ErrorWithCode(c, response.CodeUnauthorized, "token 格式错误")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析验证 token
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			response.ErrorWithCode(c, response.CodeUnauthorized, "token 无效或已过期")
			c.Abort()
			return
		}

		// 学习要点：验证通过后，把用户信息存入 context
		// 后续 controller 中用 c.Get("user_id") 取出使用
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		// 学习要点：c.Next() 继续执行后续中间件和 controller
		c.Next()
	}
}
