package models

import "time"

// Message 对话记录
type Message struct {
	ID        uint   `gorm:"primaryKey"`
	SessionID string `gorm:"index"`
	Role      string
	Content   string `gorm:"type:text"`
	CreatedAt time.Time
}