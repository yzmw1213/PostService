package model

import (
	"time"
)

// Comment コメント構造体
type Comment struct {
	CommentID      uint32 `gorm:"primary_key"`
	PostID         uint32 `validate:"required,number"`
	CreateUserID   uint32 `validate:"required,number"`
	CommentContent string `validate:"min=1,max=120"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// JoinComment コメント紐付け構造体
type JoinComment struct {
	Comment    Comment
	CreateUser User
}
