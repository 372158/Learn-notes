package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name 		string	`json:"name"`
	Age			int		`json:"age"`
	Email		string	`json:"email"`
	Password	string	`json:""`
	Articles	[]Article	`json:"articles,omitempty"`
	Roles		[]Role		`gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

type Article struct {
	gorm.Model
	Title		string 	`json:"title"`
	Content		string	`json:"content"`
	UserID		uint	`json:"user_id"`
}

type Role struct {
	gorm.Model
	Name		string	`json:"name"`
	Description	string	`json:"description"`
	Users		[]User	`gorm:"many2many:user_roles;" json:"users,omitempty"`
}

//	知识点 1：定义 Scopes （查询作用域）

//	AgeGreaterThan	年龄大于指定值
func	AgeGreaterThan(age int)	 func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where("age > ?", age)
		}
}

// AgeLessThan 年龄小于指定值
func AgeLessThan(age int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("age < ?", age)
	}
}

//	NameContains 姓名包含指定关键词（模糊查询）
func NameContains(kw string) func(db *gorm.DB)	*gorm.DB {
	return	func(db *gorm.DB) *gorm.DB {
		if kw == "" {
			return db
		}
		return db.Where("name LIKE ?", "%" + kw + "%")
	}
}


//	CreatedAfter 船舰时间在指定日期之后
func CreatedAfter(t time.Time)	func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("created_at > ?", t)
	}
}


//	Witharticles 预加载文章（用于查询时自动带出关联数据）
func WithArticles(db *gorm.DB) *gorm.DB {
	return db.Preload("Articles")
}

//	WithRoles 预加载角色
func WithRoles(db *gorm.DB) *gorm.DB {
	return	db.Preload("Roles")
}


//	Paginate 分页（page 从 1 开始）
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page < 1 {
			page = 1
		}
		if pageSize < 1 {
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}


// OrderByCreatedDesc	按创建时间降序排列
func OrderByCreatedDesc(db *gorm.DB) *gorm.DB {
	return db.Order("created_at DESC")
}


func main() {
	// 连接数据库
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数库失败：" + err.Error())
	}
	db.AutoMigrate(&User{}, &Article{}, &Role{})


	//初始化一些测试数据
	initTestData(db)

	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		//          用户查询（展示 Scopes 的威力）
		v1.GET("/users", func(c *gin.Context) {
			// 从 Query 参数获取筛选条件
			minAge, _ := strconv.Atoi(c.DefaultQuery("min_age", "0"))
			maxAge, _ := strconv.Atoi(c.DefaultQuery("max_age", "99"))

			keyword := c.Query("keyword")
			page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
			pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

			//是否预加载关联
			withArticles := c.Query("with_articles") == "true"
			withRoles := c.Query("with_roles") == "true"


			var users []User
			var total int64


			//			知识点 3 ：链式调用 Scopes
			//			先构建基础查询（用于计数）
			query := db.Model(&User{})
			//	应用筛选条件
			query = query.Scopes(
				AgeGreaterThan(minAge),
				AgeLessThan(maxAge),
				NameContains(keyword),
				OrderByCreatedDesc,
			)

			//	先计数（总条数）
			query.Count(&total)

			//	再查数据（应用分页 + 预加载 ）
			dataQuery := query
			//预加载文章
			if withArticles {
				dataQuery = dataQuery.Scopes(WithArticles)
			}
			//预加载角色
			if withRoles {
				dataQuery = dataQuery.Scopes(WithRoles)
			}
			//分页
			dataQuery = dataQuery.Scopes(Paginate(page, pageSize))

			dataQuery.Find(&users)

			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"data": gin.H{
					"list": users,
					"total": total,
					"page":	page,
					"size":	pageSize,
				},
			})
		})

		//			知识点 4: 预定义常用查询（进一步封装）
		//	查询成年用户（年龄 >= 18）
		v1.GET("/users/adults", func(c *gin.Context) {
			var users []User
			// 可以直接祝贺多个 Scopes
			db.Scopes(
				AgeGreaterThan(17),	//	年龄大于 17 即 >= 18
				WithArticles,
				OrderByCreatedDesc,
			).Find(&users)
			c.JSON(http.StatusOK, gin.H{"code": 0, "data": users})
		})

		//		查询未成年用户（年龄 < 18）
		v1.GET("/users/minors", func(c *gin.Context) {
			var users []User
			db.Scopes(
				AgeLessThan(18),
				WithArticles,
			).Find(&users)
			c.JSON(http.StatusOK, gin.H{"code": 0, "data": users})
		})

		// 		其他基础接口（简化版）
		v1.POST("/users", func(c *gin.Context) {
			var user User
			c.ShouldBindJSON(&user)
			db.Create(&user)
			c.JSON(http.StatusCreated, gin.H{"code": 0, "data": user})
		})

		v1.GET("/user/:id", func(c *gin.Context) {
			id, _ := strconv.Atoi(c.Param("id"))
			var user User
			if err := db.Preload("Articles").Preload("Roles").First(&user, id).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "用户不存在"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": 0, "data": user})
		})
		//		文章查询（也是用Scopes）
		v1.GET("/articles", func(ctx *gin.Context) {
			var articles	[]Article
			//按标题模糊搜索
			keyword := ctx.Query("keyword")
			query := db.Model(&Article{})
			if keyword != "" {
				query = query.Where("title LIKE ?", "%" + keyword + "%")
			}
			query.Preload("User").Find(&articles)
			ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": articles})
		})

		v1.POST("/articles", func(ctx *gin.Context) {
			var article Article
			ctx.ShouldBindJSON(&article)
			db.Create(&article)
			ctx.JSON(http.StatusCreated, gin.H{"code": 0, "data": article})
		})

	}

	r.Run(":8080")
}

//		辅助函数： 初始化测试数据
func initTestData(db *gorm.DB) {
	//	检查是否已有数据， 避免重复插入
	var count int64
	db.Model(&User{}).Count(&count)
	if count > 0 {
		return
	}

	users := []User {
		{Name: "张三", Age: 18, Email: "zhangsan@example.com"},
		{Name: "李四", Age: 25, Email: "lisi@example.com"},
		{Name: "王五", Age: 30, Email: "wangwu@example.com"},
		{Name: "赵六", Age: 16, Email: "zhaoliu@example.com"},
		{Name: "孙七", Age: 22, Email: "sunqi@example.com"},
		{Name: "周八", Age: 35, Email: "zhouba@example.com"},
	}
	db.Create(&users)
}
