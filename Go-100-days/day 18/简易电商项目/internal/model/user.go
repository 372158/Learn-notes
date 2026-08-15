// internal/model/user.go
// 用户数据模型，对应数据库 users 表
// 学习要点：GORM 模型用结构体定义，结构体字段对应表字段
//
// 为什么用模型而非直接写 SQL：
// 1. 类型安全，编译期检查
// 2. AutoMigrate 自动建表，不用手写 DDL
// 3. 字段标签控制 JSON 序列化和数据库约束

package model

import "time"

// User 用户模型
// 学习要点：结构体标签（tag）的作用
// - json tag: 控制 JSON 序列化时的字段名，"-" 表示不序列化（密码不返回前端）
// - gorm tag: 控制数据库字段属性（主键、唯一索引、非空等）
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`                 // 主键，GORM 默认自增
	Username  string    `json:"username" gorm:"uniqueIndex;size:50;not null"` // 用户名，唯一索引
	Password  string    `json:"-" gorm:"size:255;not null"`           // 密码，json:"-" 表示不返回前端
	Phone     string    `json:"phone" gorm:"size:20"`                 // 手机号
	CreatedAt time.Time `json:"created_at"`                           // 创建时间，GORM 自动维护
	UpdatedAt time.Time `json:"updated_at"`                           // 更新时间，GORM 自动维护
}

// TableName 自定义表名
// 学习要点：默认 GORM 会把结构体名转为复数小写（User -> users）
// 如果想自定义表名，实现 TableName 方法即可
func (User) TableName() string {
	return "users"
}
