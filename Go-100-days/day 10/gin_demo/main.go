package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	//船舰gin引擎
	engine := gin.Default()
	engine.GET("/ping",func (context *gin.Context)  {
		context.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	//开启服务器，默认监听localhost:8080
	engine.Run()
}
