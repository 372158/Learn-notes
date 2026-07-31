package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/post",func(ctx *gin.Context) {
		id := ctx.Query("id")
		page := ctx.DefaultQuery("page", "0")
		name := ctx.PostForm("name")
		message := ctx.PostForm("message")

		fmt.Printf("id :%s; page: %s; name: %s; message: %s\n", id, page, name, message)
		ctx.String(http.StatusOK, "id: %s; page: %s; name: %s; message: %s", id, page, name, message)
		
	})

	router.Run(":8080")
}