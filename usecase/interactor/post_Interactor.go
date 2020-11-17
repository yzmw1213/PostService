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
	postTag  model.PostTag
	postTags []model.PostTag
	rows     *sql.Rows
	validate *validator.Validate
)

// PostInteractor 投稿サービスを提供するメソッド群
type PostInteractor struct{}

var _ repository.PostRepository = (*PostInteractor)(nil)

// Create 投稿1件を作成
func (b *PostInteractor) Create(postData *model.JoinPost) (*model.JoinPost, error) {
	validate = validator.New()
	post := postData.Post
	tags := postData.PostTags

	// Post構造体のバリデーション
	if err := validate.Struct(post); err != nil {
		return postData, err
	}

	// トランザクション開始
	tx := db.StartBegin()

	// 投稿登録
	if err := tx.Create(post).Error; err != nil {
		db.EndRollback()
		return postData, err
	}

	postID := post.ID
	// 投稿とタグ紐付け情報登録
	for _, tag := range tags {
		tag.PostID = postID
		if err := tx.Create(tag).Error; err != nil {
			db.EndRollback()
			return postData, err
		}
	}
	// トランザクションを終了しコミット
	db.EndCommit()
	return postData, err
}

// DeleteByID 指定されたIDに対する投稿1件を削除
func (b *PostInteractor) DeleteByID(id uint32) error {
	var post model.Post
	var postTag model.PostTag

	// トランザクション開始
	tx := db.StartBegin()
	// 指定されたPostIDのPostを削除
	if err := tx.Where("id = ?", id).Delete(&post).Error; err != nil {
		db.EndRollback()
		return err
	}
	// 指定されたPostIDのPostTagを削除
	if err := tx.Where("post_id = ?", id).Delete(&postTag).Error; err != nil {
		db.EndRollback()
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

// GetByID IDを元に投稿を1件取得する
func (b *PostInteractor) GetByID(ID uint32) (model.Post, error) {
	DB := db.GetDB()
	row := DB.First(&post, ID)
	if err := row.Error; err != nil {
		log.Printf("Error happend while Read for ID: %v\n", ID)
		return model.Post{}, err
	}
	DB.Table(db.PostTableName).Scan(row)
	return post, nil
}

// attachJoinData 投稿物の

// listPostTagsByID PostIDを元にpostTagを検索し返す
func listPostTagsByID(ID uint32) ([]model.PostTag, error) {
	var postTagList []model.PostTag
	DB := db.GetDB()
	rows, err := DB.Where("post_id = ?", ID).Find(&postTags).Rows()
	if err != nil {
		log.Println("Error occured")
		return nil, err
	}
	for rows.Next() {
		DB.ScanRows(rows, &postTag)
		postTagList = append(postTagList, postTag)
	}

	return postTagList, nil
}

func countPostTag() int {
	var count int
	DB := db.GetDB()
	DB.Model(&postTag).Count(&count)
	return count
}
