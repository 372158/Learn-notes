// internal/dao/user.go
// 用户数据访问层（DAO: Data Access Object）
// 学习要点：DAO 层只负责数据库读写，不包含业务逻辑
//
// 为什么要分 DAO 层：
// 1. 职责单一：DAO 只管数据库操作，service 管业务逻辑
// 2. 易测试：可以 mock DAO 层来测试 service
// 3. 易替换：换数据库（如 MySQL 换 PostgreSQL）只改 DAO 层
//
// 约定：DAO 层方法返回 error 而非直接 panic，由上层处理错误

package dao

import (
	"simple-ecommerce/internal/model"
	"simple-ecommerce/initialize"
)

// UserDAO 用户数据访问对象
// 学习要点：用结构体封装 DAO，便于扩展（如未来加缓存）
type UserDAO struct{}

// NewUserDAO 创建 UserDAO 实例
// 学习要点：用 New 函数创建实例，是 Go 的惯用写法
func NewUserDAO() *UserDAO {
	return &UserDAO{}
}

// CreateUser 创建用户
// 学习要点：GORM Create 方法插入数据
// db.Create(&user) 会把 user 插入数据库，并回填自增 ID
func (dao *UserDAO) CreateUser(user *model.User) error {
	return initialize.DB.Create(user).Error
}

// GetUserByUsername 根据用户名查询用户
// 学习要点：GORM 查询单条记录用 First 或 Take
// - First: 没找到返回 ErrRecordNotFound
// - Where: 条件查询
func (dao *UserDAO) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := initialize.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID 根据 ID 查询用户
func (dao *UserDAO) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := initialize.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
