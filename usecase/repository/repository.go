package repository

import "github.com/yzmw1213/PostService/domain/model"

// PostRepository 投稿サービスの抽象定義
type PostRepository interface {
	Create(*model.JoinPost) (*model.JoinPost, error)
	GetByID(id uint32) (model.Post, error)
	GetJoinPostByID(id uint32) (model.JoinPost, error)
	DeleteByID(id uint32) error
	List() ([]model.JoinPost, error)
	Update(*model.JoinPost) (*model.JoinPost, error)
	Like(*model.PostLikeUser) (*model.PostLikeUser, error)
	NotLike(*model.PostLikeUser) (*model.PostLikeUser, error)
	CreateComment(*model.Comment) (*model.Comment, error)
	UpdateComment(*model.Comment) (*model.Comment, error)
	DeleteComment(id uint32) error
}

// TagRepository タグサービスの抽象定義
type TagRepository interface {
	Create(*model.Tag) (*model.Tag, error)
	DeleteByID(uint32) error
	GetTagByTagName(string) (model.Tag, error)
	ListValidTag() ([]model.Tag, error)
	List() ([]model.Tag, error)
	Update(*model.Tag) (*model.Tag, error)
}
