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
