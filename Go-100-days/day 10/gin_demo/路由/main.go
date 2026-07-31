package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world!")
	})

	router.POST("/users", func(ctx *gin.Context) {
		name := ctx.PostForm("name")
		ctx.JSON(http.StatusCreated, gin.H{"user": name})
	})
	router.Run(":8080")
}