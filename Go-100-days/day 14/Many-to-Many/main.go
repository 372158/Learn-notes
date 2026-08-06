package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//模型定义

//User 用户模型
type User struct {
	gorm.Model
	Name 	string	`json:"name"`
	Age		int	`json:"age"`
	Articles	[]Article	`json:"articles,omitempty"` 	//一对多：一个用户拥有多篇文章
	Roles	[]Role	`gorm:"many2many:user_roles;" json:"roles,omitempty"`
	// 多对多： 一个用户有多个角色
}

//Article 文章模型
type Article struct {
	gorm.Model
	Title	string	`json:"title"`
	Content	string	`json:"content"`
	UserID	uint	`josn:"user_id"` //外键： 属于哪个用户
}

//	 Role 角色模型（新增）
type Role struct {
	gorm.Model
	Name		string `json:"name"`
	Description	string	`json:"description"`
	Users		[]User	`gorm:"many2many:user_roles;" json:"users,omitempty"`
	//多对多反向引用
}

func main()	{
	// 	连接数据库
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("链接数据库失败：" + err.Error())
	}

	db.AutoMigrate(&User{}, &Article{}, &Role{})

	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		// 用户相关

		//创建用户（带文章和角色）
		v1.POST("/users", func(c *gin.Context) {
			var input struct {
				Name	string	`json:"name"`
				Age		int	`json:"age"`
				Articles	[]struct {
					Title	string	`json:"title"`
					Content string	`json:"content"`
				}	`json:"articles"`
				RoleIDs	[]uint	`json:"role_ids"` // 前端传角色IDl列表
			}

			if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
				return
			}

			//构建用户
			user := User{Name: input.Name, Age: input.Age}

			// 构建关联文章
			for _, a := range input.Articles {
				user.Articles = append(user.Articles, Article{
					Title: a.Title,
					Content: a.Content,
				})

			}

			//  知识点： 多对多关联
			//如果前端传了角色ID 列表，我们需要把对应的角色查出来，关联到用户。
			// 注意： 此时角色必须已经存在于数据库中（先创建角色，再分配角色）

			if len(input.RoleIDs) > 0 {
				var  roles []Role
				db.Find(&roles, input.RoleIDs)	//	批量查询角色 
				user.Roles = roles				// f赋值给用户的 Rolses字段
			}

			db.Create(&user)
			c.JSON(http.StatusCreated, gin.H{"code": 0, "data": user})
		})


		//查询所有用户（预加载文章 + 角色）
		v1.GET("/users", func(c *gin.Context) {
			var users []User
			//知识点： 多个 Preload 可以链式调用

			db.Preload("Artticlse").Preload("Roles").Find(&users)
			c.JSON(http.StatusOK, gin.H{"code": 0, "data": users})
		})


		//查询单个用户
		v1.GET("/uers/:id", func(c *gin.Context) {
			id, _ := strconv.Atoi(c.Param("id"))
			var user User
			if err := db.Preload("Articles").Preload("Roles").First(&user, id).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "用户不存在"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": 0, "data": user})
		})


		//给用户分配角色（追加）
		v1.POST("/users/:id/roles", func(c *gin.Context) {
			id, _ := strconv.Atoi(c.Param("id"))
			var input struct {
				RolesIDs	[]uint	`json:"roles_ids"`
			}
			c.ShouldBindJSON(&input)

			var user User
			if err := db.First(&user, id).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "用户不存在"})
				return
			}

			//查询要添加的角色
			var roles []Role
			db.Find(&roles, input.RolesIDs)

			//	知识点： append 追加关联
			//	Model(&user).Association("Roles").Append(roles) 的作用

			//在中间表 user_roles 中插入心得关联记录， 不会删除已有的关联。
			//相当于给用户添加新角色， 保留旧角色。
			db.Model(&user).Association("Roles").Append(roles)

			//重新加载用户数据（带上角色）
			db.Preload("Roles").First(&user, id)
			c.JSON(http.StatusOK, gin.H{"code": 0, "data": user})
		})


		//移除用户的某个角色
		v1.DELETE("/users/:id/roles/:roleId", func(c *gin.Context) {
			userID, _ := strconv.Atoi(c.Param("id"))
			roleID, _ := strconv.Atoi(c.Param("roleId"))

			var  user User
			db.First(&user, userID)
			var role Role
			db.First(&role, roleID)

			// 知识点： Delete 删除关联
			//	从中间表 user_roles中删除这条关联记录， 但不会删除角色本身
			//	角色依然存在，只是这个用户不在拥有这个角色。
			db.Model(&user).Association("Roles").Delete(&role)

			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "移除角色成功"})
		})


		//		角色相关
		//		创建角色
		v1.POST("/roles", func (c *gin.Context)  {
			var role Role
			c.ShouldBindJSON(&role)
			db.Create(&role)
			c.JSON(http.StatusCreated, gin.H{"code": 0, "data": role})

		})

		v1.GET("/roles", func(c *gin.Context) {
			var roles []Role
			db.Preload("Users").Find(&roles)
			c.JSON(http.StatusOK, gin.H{"code": 0, "data": roles})
		})
	}

	r.Run(":8080")
}