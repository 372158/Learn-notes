package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//User 结构体
type User struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Age int `json:"age"`
}

//模拟数据库（内存切片）
var users = []User{
	{ID: 1, Name:"张三", Age: 18},
	{ID: 2, Name:"李四", Age: 20},
}

var nextID = 3

func main() {
	r := gin.Default()

	//路由分组
	v1 := r.Group("/api/v1")
	{
		//查询所有用户
		v1.GET("/users", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"data": users,
			})
		})

		//查询单个用户
		v1.GET("/users/:id", func(c *gin.Context) {
			id, _ := strconv.Atoi(c.Param("id"))
			for _, u := range users {
				if u.ID == id {
					c.JSON(http.StatusOK, gin.H{
						"code": 0,
						"data": u,
					})
					return
				}
			}
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"message": "用户不存在",
			})
		})

		//创建用户
		v1.POST("/users", func(c *gin.Context) {
			var newUser User
			if err := c.ShouldBindBodyWithJSON(&newUser); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code": 404,
					"message": "参数错误",
				})
				return
			}
			newUser.ID = nextID
			nextID++
			users = append(users, newUser)
			c.JSON(http.StatusCreated, gin.H{
				"code": 0,
				"data": newUser,
			})
		})

		//更新用户
		v1.PUT("/users/:id",func (c *gin.Context)  {
			id, _ := strconv.Atoi(c.Param("id"))
			var updateData User
			if err := c.ShouldBindJSON(&updateData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code": 404,
					"message": "参数错误",
				})
				return
			}

			for i, u := range users {
				if u.ID == id {
					users[i].Name = updateData.Name
					users[i].Age = updateData.Age
					c.JSON(http.StatusOK, gin.H{
						"code": 0,
						"data": users[i],
					})
					return
				}
			}

			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"message": "用户不存在",
			})
		})

		//删除用户
		v1.DELETE("/users/:id", func (c *gin.Context)  {
			id, _ := strconv.Atoi(c.Param("id"))
			for i, u := range users {
				if u.ID == id {
					//切片中删除
					users = append(users[:i], users[i+1:]...)
					c.JSON(http.StatusOK, gin.H{
						"code": 0,
						"message": "删除成功",
					})
					return
				}
			}
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"message": "用户不存在",
			})
			
		})

		r.Run(":8080")
	}
}