package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/worl", func (c *gin.Context)  {
		ctx := c.Request.Context()

		for i := 1; i <= 10; i++ {
			select {
			case <- time.After(1 *time.Second):
				fmt.Printf("[服务端] 干活中 %d/10\n", i)

			case <- ctx.Done():
				fmt.Println("[服务端] 客户端断开，停止干活")
				return
			}
		}
		c.String(200, "10 秒的活干完了")
	})

	r.Run(":8080")
}