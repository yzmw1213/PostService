package repository

import "github.com/yzmw1213/PostService/domain/model"

// PostRepository 投稿サービスの抽象定義
type PostRepository interface {
	Create(*model.Post) (*model.Post, error)
	Delete(*model.Post) error
	List() ([]model.Post, error)
	Update(*model.Post) (*model.Post, error)
}
