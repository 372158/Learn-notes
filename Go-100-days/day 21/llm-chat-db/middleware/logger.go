package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 请求日志中间件：记录每个请求的方法、路径、状态码、耗时
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next() // 执行后续 handler（含真正处理请求的回调）

		// 处理完回来，记录日志
		log.Printf("[请求] %s %s 状态=%d 耗时=%v",
			c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(start))
	}
}