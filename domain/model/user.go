package model

type User struct {
	ID        uint32 `gorm:"primary_key"`
	UserName  string `validate:"min=6,max=16"`
	Gender    uint32 `validate:"oneof=0 1 2 9"`
	Authority uint32 `validate:"oneof=0 1 2 9"`
}
