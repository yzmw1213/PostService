package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// Post 投稿サービス構造体
type Post struct {
	ID        uint32 `gorm:"primary_key"`
	UserID    uint32 `validate:"required,number"`
	Content   string `validate:"min=1,max=32"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
