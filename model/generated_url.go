package model

import (
	"time"
)

type GeneratedUrl struct {
	UrlID       int64     `gorm:"primaryKey;column:url_id"`
	Username    string    `gorm:"column:username"`
	OriginalUrl string    `gorm:"column:original_url"`
	TinyUrl     string    `gorm:"column:tiny_url;unique"`
	IsActive    bool      `gorm:"column:is_active;default:true"`
	StartTime   time.Time `gorm:"column:start_time"`
	EndTime     time.Time `gorm:"column:end_time"`
}

// TableName specifies the table name for the GeneratedUrl model
func (GeneratedUrl) TableName() string {
	return "generatedurl"
}
