package model

import "time"

// Tag タグサービス構造体
type Tag struct {
	ID           int32  `gorm:"primary_key"`
	TagName      string `validate:"required"`
	CreateUserID string `validate:"required,alphanum"`
	UpdateUserID string
	Status       int32 `validate:"number"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
