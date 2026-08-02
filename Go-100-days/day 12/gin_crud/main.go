package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string `json:"name"`
	Age int `json:"age"`
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败：" + err.Error())
	}
	db.AutoMigrate(&User{})

	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.GET("/users", func(c *gin.Context) {
			var users []User
			db.Find(&users)
			c.JSON(http.StatusOK, gin.H{"code": 200, "data": users})
		})

		v1.GET("/users/:id", func(c *gin.Context) {
			id, _ := strconv.Atoi(c.Param("id"))
			var user User
			if err := db.First(&user, id).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "用户不存在"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": 200, "data": user})

		})

		v1.POST("/users", func(c *gin.Context) {
			var user User
			if err := c.ShouldBindJSON(&user); err != nil {
				c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "参数错误"})
				return
			}
			db.Create(&user)
			c.JSON(http.StatusOK, gin.H{"code": 200, "data": user})
		})

		v1.PUT("/users/:id", func(c *gin.Context) {
			id, _ := strconv.Atoi(c.Param("id"))
			var user User
			if err := db.First(&user, id).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"code": 404,"msg": "用户不存在"})
				return
			}
			var updateData User
			c.ShouldBindJSON(&updateData)
			db.Model(&user).Updates(User{Name: updateData.Name, Age: updateData.Age})
			c.JSON(http.StatusOK, gin.H{"code": 0, "data": user})
		})	

		v1.DELETE("/users/:id", func (c *gin.Context)  {
			id, _ := strconv.Atoi(c.Param("id"))
			result := db.Delete(&User{}, id)
			if result.RowsAffected == 0{
				c.JSON(http.StatusNotFound, gin.H{"code": 404,"msg": "用户不存在"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除成功"})
		})

		r.Run(":8080")
	}
}