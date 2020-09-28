package db

import (
	"fmt"
	"os"

	// gormのmysql接続用
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yzmw1213/PostService/domain/model"

	"github.com/jinzhu/gorm"
)

var (
	// DB データベース構造体
	DB *gorm.DB
	// PostTableName 投稿サービステーブル名
	PostTableName string = "posts"
	// TagTableName タグサービステーブル名
	TagTableName string = "tags"
)

func initDB() {
	var err error
	DBMS := "mysql"
	DBNAME := os.Getenv("DB_NAME")
	PASSWORD := os.Getenv("DB_PASSWORD")
	USER := os.Getenv("DB_USER")
	PROTOCOL := fmt.Sprintf("tcp(%s)", os.Getenv("DB_ADRESS"))
	OPTION := "?charset=utf8mb4&parseTime=True&loc=Local"
	CONNECTION := fmt.Sprintf("%s:%s@%s/%s%s", USER, PASSWORD, PROTOCOL, DBNAME, OPTION)

	DB, err = gorm.Open(DBMS, CONNECTION)
	if err != nil {
		panic(err)
	}
}

// Init DB接続と、マイグレーションを行う。
func Init() {
	initDB()
	// マイグレーション実行
	autoMigration()
}

// Close DBと切断する。
func Close() {
	if err := DB.Close(); err != nil {
		panic(err)
	}
}

// GetDB DB接続情報を返す
func GetDB() *gorm.DB {
	initDB()
	return DB
}

func autoMigration() {
	fmt.Println("migration")
	err := DB.AutoMigrate(&model.Post{}).Error
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&model.Tag{}).Error
	if err != nil {
		panic(err)
	}
}
