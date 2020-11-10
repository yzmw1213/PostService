package interactor

import (
	"context"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/yzmw1213/PostService/db"
	"github.com/yzmw1213/PostService/domain/model"
	"github.com/yzmw1213/PostService/usecase/repository"
)

var (
	tag  model.Tag
	tags []model.Tag
	// ValidTagStatus タグ公開ステータス
	ValidTagStatus uint32 = 1
	// InValidTagStatus タグ非公開ステータス
	InValidTagStatus uint32 = 2
)

// TagInteractor タグサービスを提供するメソッド群
type TagInteractor struct{}

var _ repository.TagRepository = (*TagInteractor)(nil)

// Create タグ1件を作成
func (i *TagInteractor) Create(postData *model.Tag) (*model.Tag, error) {
	validate = validator.New()
	DB := db.GetDB()

	// Tag構造体のバリデーション
	if err := validate.Struct(postData); err != nil {
		return postData, err
	}
	if err := DB.Create(postData).Error; err != nil {
		return postData, err
	}

	return postData, err
}

// DeleteByID 指定したIDのタグ1件を削除
func (i *TagInteractor) DeleteByID(id uint32) error {
	DB := db.GetDB()
	if err := DB.Where("id = ? ", id).Delete(&tag).Error; err != nil {
		return err
	}
	return nil
}

// List タグを全件取得
func (i *TagInteractor) List() ([]model.Tag, error) {
	var tagList []model.Tag
	rows, err := listAllTag(context.Background())
	if err != nil {
		fmt.Println("Error happened")
		return []model.Tag{}, err
	}
	for _, row := range rows {
		tagList = append(tagList, row)
	}

	return tagList, nil
}

// ListValidTag 有効タグを全件取得する
func (i *TagInteractor) ListValidTag() ([]model.Tag, error) {
	DB := db.GetDB()
	var tags []model.Tag

	if err := DB.Where("status = ?", "1").Find(&tags).Error; err != nil {
		return []model.Tag{}, err
	}
	return tags, nil
}

// listAllTag タグ全件取得
func listAllTag(ctx context.Context) ([]model.Tag, error) {
	DB := db.GetDB()

	_ = DB.Find(&tags)
	if err != nil {
		log.Println("Error occured")
		return nil, err
	}
	return tags, nil
}

// ListAllValidTag 有効タグ全件取得
func (i *TagInteractor) ListAllValidTag() ([]model.Tag, error) {
	DB := db.GetDB()
	var tags []model.Tag
	err := DB.Where("status = ?", ValidTagStatus).Select("tags.id, tags.tag_name").Find(&tags).Error
	log.Println("tags", tags)
	if err != nil {
		fmt.Println("Error happened")
		return []model.Tag{}, err
	}

	return tags, nil
}

// Update タグを更新する
func (i *TagInteractor) Update(postData *model.Tag) (*model.Tag, error) {
	DB := db.GetDB()
	validate = validator.New()

	// Tag構造体のバリデーション
	if err := validate.Struct(postData); err != nil {
		return postData, err
	}
	if err := DB.Model(&tag).Update(&postData).Error; err != nil {
		return postData, err
	}

	return postData, nil
}

// GetTagByTagName TagNameを元にタグを1件取得する
func (i *TagInteractor) GetTagByTagName(tagName string) (model.Tag, error) {
	var tag model.Tag

	DB := db.GetDB()
	row := DB.Where("tag_name = ?", tagName).First(&tag)
	if err := row.Error; err != nil {
		return tag, err
	}
	DB.Table(db.TagTableName).Scan(row)

	return tag, nil
}

// GetTagByTagID TagIDを元にタグを1件取得する
func (i *TagInteractor) GetTagByTagID(tagID uint32) (model.Tag, error) {
	var tag model.Tag

	DB := db.GetDB()
	row := DB.Where("id = ?", tagID).First(&tag)
	if err := row.Error; err != nil {
		return tag, err
	}
	DB.Table(db.TagTableName).Scan(row)

	return tag, nil
}
