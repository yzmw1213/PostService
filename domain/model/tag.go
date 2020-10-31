package model

import "time"

// Tag タグサービス構造体
type Tag struct {
	ID           uint32 `gorm:"primary_key"`
	TagName      string `validate:"min=1,max=12"`
	CreateUserID uint32 `validate:"required,number"`
	UpdateUserID uint32 `validate:"number"`
	Status       uint32 `validate:"number"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
