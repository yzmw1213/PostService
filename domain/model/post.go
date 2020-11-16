package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// Post 投稿サービス構造体
type Post struct {
	ID uint32 `gorm:"primary_key"`
	// Status       uint32 `validate:"required,number"`
	Title        string `validate:"min=1,max=32"`
	Content      string `validate:"min=1,max=240"`
	MaxNum       uint32 `validate:"required,number"`
	Gender       uint32 `validate:"oneof=1 2 3"`
	CreateUserID uint32 `validate:"required,number"`
	UpdateUserID uint32 `validate:"number"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
