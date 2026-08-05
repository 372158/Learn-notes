package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 1. 模型定义（替换之前的 User 结构体）
type User struct {
	gorm.Model          // 自带 ID, CreatedAt, UpdatedAt, DeletedAt
	Name         string `json:"name"`
	Age          int    `json:"age"`
}

func main() {
	// 2. 连接数据库（替换之前的 var users = []User{}）
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败: " + err.Error())
	}

	// 3. 自动建表（不用手动写 SQL）
	db.AutoMigrate(&User{})

	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		// -------- 查询所有（替换 for 循环） --------
		v1.GET("/users", func(c *gin.Context) {
			var users []User
			var total int64
			query := db.Model(&User{})

			//1.年龄范围查询
			minAge := c.Query("minAge")
			maxAge := c.Query("maxAge")
			
			if minAge != "" {
				query = query.Where("age >= ?", minAge)
			}
			
			if maxAge != "" {
				query = query.Where("age <= ?", maxAge)
			}

			//2.关键字模糊查询
			keyword := c.Query("keyword")
			if keyword != "" {
				query = query.Where("name LIKE ?", "%" + keyword +"%")
			}

			//3. 分页参数（默认第1页， 每页10条)
			page, _ := strconv.Atoi(c.DefaultQuery("page","1"))
			size, _ := strconv.Atoi(c.DefaultQuery("size","10"))
			offset := (page - 1) * size

			//4. 执行查询（先计数，在查数据）
			query.Count(&total)
			query.Offset(offset).Limit(size).Find(&users)

			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"data": gin.H{
					"list": users,
					"total": total,
					"page": page,
					"size": size,
				},
			})
		})

		// -------- 查询单个（替换 for 遍历查找） --------
		v1.GET("/users/:id", func(c *gin.Context) {
			id, _ := strconv.Atoi(c.Param("id"))
			var user User
			// 等价于 SELECT * FROM users WHERE id = ? LIMIT 1
			if err := db.First(&user, id).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "用户不存在"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": 0, "data": user})
		})

		// -------- 创建用户（替换 append） --------
		v1.POST("/users", func(c *gin.Context) {
			var newUser User
			if err := c.ShouldBindJSON(&newUser); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
				return
			}
			db.Create(&newUser) // 等价于 INSERT INTO users (name, age) VALUES (?, ?)
			c.JSON(http.StatusCreated, gin.H{"code": 0, "data": newUser})
		})


		//批量创建用户（事务）
		v1.POST("/users/batch", func(c *gin.Context) {
			var users []User
			if err := c.ShouldBindJSON(&users); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数格式错误，请传入JSON数组"},)
				return
			}

			//2. 开启事务执行批量创建
		//注意：这里使用db.Transaction, 如果返回error 则自动回滚，返回nil 则自动提交
		err := db.Transaction(func(tx *gorm.DB) error {
			//批量插入（GORM 会自动拆分成一条SQL 或 分批插入）
			if err := tx.Create(&users).Error; err != nil {
				//只要插入报错（比如字段超长），就返回错误，出发回滚
				return err
			}
			//也可以在这里加其他表的操作，保持原子性
			return nil
		})

		//3.处理事务结果
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg": "批量创建失败，事务已回滚",
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"code": 0,
			"data": users,//返回创建成功后的数据，（ID 和时间戳会被自动填回）
		})
		})


		

		// -------- 更新用户（替换切片索引修改） --------
		v1.PUT("/users/:id", func(c *gin.Context) {
			id, _ := strconv.Atoi(c.Param("id"))
			var user User
			if err := db.First(&user, id).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "用户不存在"})
				return
			}

			var updateData User
			c.ShouldBindJSON(&updateData)

			// 等价于 UPDATE users SET name=?, age=? WHERE id = ?
			db.Model(&user).Updates(User{Name: updateData.Name, Age: updateData.Age})
			c.JSON(http.StatusOK, gin.H{"code": 0, "data": user})
		})

		// -------- 删除用户（替换 append 拼接删除） --------
		v1.DELETE("/users/:id", func(c *gin.Context) {
			id, _ := strconv.Atoi(c.Param("id"))
			// 等价于 DELETE FROM users WHERE id = ? （GORM 默认是软删除）
			result := db.Delete(&User{}, id)
			if result.RowsAffected == 0 {
				c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "用户不存在"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除成功"})
		})
	}
	r.Run(":8080")
}