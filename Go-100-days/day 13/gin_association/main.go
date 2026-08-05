package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User 用户模型（一对多的 "一"）
type User struct {
	gorm.Model
	Name 	string	`json:"name"`
	Age		int 	`json:"age"`
	Articles []Article	`json:"articles,omitempty"`//关联的文章列表
}

//Artices  文章模型（一对多的“多”）
type Article struct {
	gorm.Model
	Title		string	`json:"title"`
	Content		string	`json:"content"`
	UserID		uint	`json:"user_id"`	//外键，指向User 的 ID
	User    User   `json:"user,omitempty"` 
}

//这是练习的第十三天
func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("链接数据库失败：" + err.Error())
	} 
	db.AutoMigrate(&User{}, &Article{})

	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		//1.创建用户（同时创建文章）
		v1.POST("/users", func(c *gin.Context) {
			var input struct {
				Name	string	`json:"title"`
				Age 	int 	`json:"age"`
				Articles	[]struct {
							Title 	string	`json:"title"`
							Content	string	`json:"content"`
				} `json:"articles"`
			}
			if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"}) 
				return
			}
			// 2.构建User 对象
			user := User{
				Name: input.Name,
				Age:  input.Age,
			}

			//3.构建关联的 Article 列表
			for _, a := range input.Articles {
				user.Articles = append(user.Articles, Article{
					Title: 	a.Title,
					Content: a.Content,	
				})
			}
			//使用十五创建用户和文章（GORM 自动处理外键)
			db.Create(&user)
			c.JSON(http.StatusCreated, gin.H{"code": 0, "data": user}) 
		})

		// 查询所有用户（预加载文章）
		v1.GET("/users", func(c *gin.Context) {
			var users []User
			//Preload("Articles") 会额外执行一条 SELECT * FROM articles WHERE user_id IN (?)
			db.Preload("Articles").Find(&users)
			c.JSON(http.StatusOK, gin.H{"code": 0, "data": users})
		})

		//查询单个用户（预加载文章）
		v1.GET("/users/:id", func(c *gin.Context) {
			id, _ := strconv.Atoi(c.Param("id"))
			var user User
			if err := db.Preload("Articles").First(&user, id).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg":"用户不存在"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": 0, "data": user})
		})

		//查询所有文章（带作者信息）
		v1.GET("/articles", func(c *gin.Context) {
			var articles []Article
			//预加载关联的 User
			db.Preload("User").Find(&articles)
			c.JSON(http.StatusOK, gin.H{"code": 0, "data": articles})
		})
	}
	r.Run(":8080")	
}