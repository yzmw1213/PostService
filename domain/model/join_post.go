package model

// JoinPost 投稿情報の紐付け構造体
type JoinPost struct {
	Post *Post
	// 紐付けられたタグ情報
	PostTags []PostTag
}
