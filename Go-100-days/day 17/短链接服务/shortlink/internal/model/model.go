package model

import (
	"time"

	"gorm.io/gorm"
)

// ShortLink 短链接模型

type ShortLink struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	ShortCode string `gorm:"uniqueIndex;size:10" json:"short_code"`
	//短码
	OriginalURL string `gorm:"type:text;not null" json:"original_url"`
	//	原始 URL
	ClickCount int `gorm:"default:0" json:"click_count"`
	//	点击次数
	ExpiredAt *time.Time `json:"expired_at,omitempty"`
	//	过期时间
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
