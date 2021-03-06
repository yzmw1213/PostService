package model

// User ユーザー構造体
type User struct {
	ID        uint32 `gorm:"primary_key"`
	UserName  string `validate:"min=6,max=16"`
	Authority uint32 `validate:"oneof=0 1 2 9"`
}
