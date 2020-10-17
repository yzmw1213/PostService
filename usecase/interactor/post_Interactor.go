package interactor

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"

	"github.com/yzmw1213/PostService/db"
	"github.com/yzmw1213/PostService/domain/model"
	"github.com/yzmw1213/PostService/usecase/repository"
)

var (
	err      error
	post     model.Post
	posts    []model.Post
	rows     *sql.Rows
	validate *validator.Validate
)

// PostInteractor 投稿サービスを提供するメソッド群
type PostInteractor struct{}

var _ repository.PostRepository = (*PostInteractor)(nil)

// Create 投稿1件を作成
func (b *PostInteractor) Create(postData *model.Post) (*model.Post, error) {
	validate = validator.New()
	DB := db.GetDB()

	// Post構造体のバリデーション
	if err := validate.Struct(postData); err != nil {
		return postData, err
	}
	if err := DB.Create(postData).Error; err != nil {
		return postData, err
	}

	return postData, err
}

// Delete 投稿1件を削除
func (b *PostInteractor) Delete(postData *model.Post) error {
	DB := db.GetDB()
	if err := DB.Delete(postData).Error; err != nil {
		return err
	}
	return nil
}

// List 投稿を全件取得
func (b *PostInteractor) List() ([]model.Post, error) {
	var postList []model.Post
	rows, err := listAll(context.Background())
	if err != nil {
		fmt.Println("Error happened")
		return []model.Post{}, err
	}
	for _, row := range rows {
		postList = append(postList, row)
	}

	return postList, nil
}

// listAll 全件取得
func listAll(ctx context.Context) ([]model.Post, error) {
	DB := db.GetDB()

	rows, err := DB.Find(&posts).Rows()
	if err != nil {
		log.Println("Error occured")
		return nil, err
	}

	for rows.Next() {
		DB.ScanRows(rows, &post)
		posts = append(posts, post)
	}
	return posts, nil
}

// Update 投稿を更新する
func (b *PostInteractor) Update(postData *model.Post) (*model.Post, error) {
	DB := db.GetDB()

	// Post構造体のバリデーション
	if err := validate.Struct(postData); err != nil {
		return postData, err
	}
	if err := DB.Model(&post).Update(&postData).Error; err != nil {
		return postData, err
	}

	return postData, nil
}

// Read IDを元に投稿を1件取得する
func (b *PostInteractor) Read(ID int32) (model.Post, error) {
	DB := db.GetDB()
	row := DB.First(&post, ID)
	if err := row.Error; err != nil {
		log.Printf("Error happend while Read for ID: %v\n", ID)
		return model.Post{}, err
	}
	DB.Table(db.PostTableName).Scan(row)
	return post, nil
}
