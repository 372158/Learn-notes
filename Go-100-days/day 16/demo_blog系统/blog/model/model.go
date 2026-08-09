package model

import (
	"gorm.io/gorm"
)
//User	用户模型
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;size:50" json:"username"`
	Password string `json:"-"`
	Email    string `gorm:"uniqueIndex;size:100" json:"email"`
	Role     string `gorm:"default:'user';size:20" json:"role"`
}


// Article 文章模型（新增）
type Article struct {
	gorm.Model
	Title		string		`gorm:"size:200" json:"title" binding:"required"`
	Content		string		`gorm:"type:text" json:"content" binding:"required"`
	Summary		string		`gorm:"size500" json:"sunmmary"`
	Views		int			`gorm:"default:0" json:"views"`
	UserID		uint		`json:"user_id"`
	User		User		`json:"user,omitempty"`
}

