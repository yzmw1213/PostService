package repository

import "github.com/yzmw1213/PostService/domain/model"

// PostRepository 投稿サービスの抽象定義
type PostRepository interface {
	Create(*model.Post) (*model.Post, error)
	Delete(*model.Post) error
	List() ([]model.Post, error)
	Update(*model.Post) (*model.Post, error)
}

// TagRepository タグサービスの抽象定義
type TagRepository interface {
	Create(*model.Tag) (*model.Tag, error)
	DeleteByID(uint32) error
	GetTagByTagName(string) (model.Tag, error)
	GetValidTag() ([]model.Tag, error)
	List() ([]model.Tag, error)
	Update(*model.Tag) (*model.Tag, error)
}
