package grpc

import (
	"context"
	"log"

	"github.com/yzmw1213/PostService/domain/model"
	"github.com/yzmw1213/PostService/grpc/postservice"
)

const (
	// StatusCreatePostSuccess 投稿作成成功ステータス
	StatusCreatePostSuccess string = "POST_CREATE_SUCCESS"
	// StatusUpdatePostSuccess 投稿更新成功ステータス
	StatusUpdatePostSuccess string = "POST_UPDATE_SUCCESS"
	// StatusLikePostSuccess 投稿お気に入り成功ステータス
	StatusLikePostSuccess string = "POST_LIKE_SUCCESS"
	// StatusNotLikePostSuccess 投稿お気に入り取り消し成功ステータス
	StatusNotLikePostSuccess string = "POST_NOTLIKE_SUCCESS"
	// StatusDeletePostSuccess 投稿削除成功ステータス
	StatusDeletePostSuccess string = "POST_DELETE_SUCCESS"
	// StatusPostNotExists 指定した投稿の登録がない時のエラーステータス
	StatusPostNotExists string = "POST_NOT_EXISTS_ERROR"
	// StatusPostTitleStringCount 件名文字数が無効のエラーステータス
	StatusPostTitleStringCount string = "POST_TITLE_COUNT_ERROR"
	// StatusPostContentStringCount 投稿内容文字数が無効のエラーステータス
	StatusPostContentStringCount string = "POST_CONTENT_COUNT_ERROR"
)

func (s server) CreatePost(ctx context.Context, req *postservice.CreatePostRequest) (*postservice.CreatePostResponse, error) {
	postData := req.GetPost()

	post := makePostModel(postData)
	tags := makePostTagModel(postData)

	joinPost := &model.JoinPost{
		Post:     post,
		PostTags: tags,
	}

	// post, tagsをJoinしてinteractor.Createに渡す
	joinPost, err := s.PostUsecase.Create(joinPost)
	if err != nil {
		return nil, err
	}

	return s.makeCreatePostResponse(StatusCreatePostSuccess), nil
}

func (s server) DeletePost(ctx context.Context, req *postservice.DeletePostRequest) (*postservice.DeletePostResponse, error) {
	id := req.GetId()

	if err := s.PostUsecase.DeleteByID(id); err != nil {
		return nil, err
	}
	return s.makeDeletePostResponse(StatusDeletePostSuccess), nil
}

func (s server) ListPost(ctx context.Context, req *postservice.ListPostRequest) (*postservice.ListPostResponse, error) {
	rows, err := s.PostUsecase.List()
	if err != nil {
		return nil, err
	}
	var posts []*postservice.Post
	for _, post := range rows {
		post := makeGrpcPost(&post)
		posts = append(posts, post)
	}
	res := &postservice.ListPostResponse{
		Post: posts,
	}
	return res, nil
}

func (s server) ReadPost(ctx context.Context, req *postservice.ReadPostRequest) (*postservice.ReadPostResponse, error) {
	ID := req.GetId()
	row, err := s.PostUsecase.GetJoinPostByID(ID)
	if err != nil {
		return nil, err
	}
	post := makeGrpcPost(&row)
	res := &postservice.ReadPostResponse{
		Post: post,
	}
	return res, nil
}

func (s server) UpdatePost(ctx context.Context, req *postservice.UpdatePostRequest) (*postservice.UpdatePostResponse, error) {
	postData := req.GetPost()

	joinPost := &model.JoinPost{
		Post:     makePostModel(postData),
		PostTags: makePostTagModel(postData),
	}
	if _, err := s.PostUsecase.Update(joinPost); err != nil {
		return nil, err
	}

	return s.makeUpdatePostResponse(StatusUpdatePostSuccess), nil
}

func (s server) LikePost(ctx context.Context, req *postservice.LikePostRequest) (*postservice.LikePostResponse, error) {
	log.Println("LikePost")
	log.Println("user", req.GetUserId())
	log.Println("post", req.GetId())
	postLikeUser := &model.PostLikeUser{
		PostID: req.GetId(),
		UserID: req.GetUserId(),
	}

	if _, err := s.PostUsecase.Like(postLikeUser); err != nil {
		return nil, err
	}
	return s.makeLikePostResponse(StatusLikePostSuccess), nil
}

func (s server) NotLikePost(ctx context.Context, req *postservice.NotLikePostRequest) (*postservice.NotLikePostResponse, error) {
	log.Println("NotLikePost")
	log.Println("user", req.GetUserId())
	log.Println("post", req.GetId())
	postLikeUser := &model.PostLikeUser{
		PostID: req.GetId(),
		UserID: req.GetUserId(),
	}

	if _, err := s.PostUsecase.NotLike(postLikeUser); err != nil {
		return nil, err
	}
	return s.makeNotLikePostResponse(StatusNotLikePostSuccess), nil
}

func makePostModel(gPost *postservice.Post) *model.Post {
	post := &model.Post{
		ID: gPost.GetId(),
		// Status:       gPost.GetStatus(),
		Title:        gPost.GetTitle(),
		Content:      gPost.GetContent(),
		CreateUserID: gPost.GetCreateUserId(),
		UpdateUserID: gPost.GetUpdateUserId(),
	}
	return post
}

func makePostTagModel(gPost *postservice.Post) []model.PostTag {
	var postTags []model.PostTag

	for _, tagID := range gPost.Tags {
		postTags = append(postTags, model.PostTag{
			PostID: gPost.Id,
			TagID:  tagID,
		})
	}
	log.Println("postTags", postTags)
	return postTags
}

func makeGrpcPost(post *model.JoinPost) *postservice.Post {
	var tags []uint32
	var likeUsers []uint32
	gPost := &postservice.Post{
		Id: post.Post.ID,
		// Status:       post.Status,
		Title:          post.Post.Title,
		Content:        post.Post.Content,
		CreateUserId:   post.Post.CreateUserID,
		CreateUserName: post.User.UserName,
		UpdateUserId:   post.Post.UpdateUserID,
	}
	// タグ
	for _, postTag := range post.PostTags {
		tags = append(tags, postTag.TagID)
	}
	gPost.Tags = tags

	// お気に入りユーザー
	for _, user := range post.PostLikeUsers {
		likeUsers = append(likeUsers, user.ID)
	}
	gPost.LikeUsers = likeUsers

	return gPost
}

func createPostRequest(post *postservice.Post) *postservice.CreatePostRequest {
	return &postservice.CreatePostRequest{
		Post: post,
	}
}

func updatePostRequest(post *postservice.Post) *postservice.UpdatePostRequest {
	return &postservice.UpdatePostRequest{
		Post: post,
	}
}

// makeCreatePostResponse CreatePostメソッドのresponseを生成し返す
func (s server) makeCreatePostResponse(statusCode string) *postservice.CreatePostResponse {
	res := &postservice.CreatePostResponse{}
	if statusCode != "" {
		responseStatus := &postservice.ResponseStatus{
			Code: statusCode,
		}
		res.Status = responseStatus
	}
	return res
}

// makeUpdatePostResponse UpdatePostメソッドのresponseを生成し返す
func (s server) makeUpdatePostResponse(statusCode string) *postservice.UpdatePostResponse {
	res := &postservice.UpdatePostResponse{}
	if statusCode != "" {
		responseStatus := &postservice.ResponseStatus{
			Code: statusCode,
		}
		res.Status = responseStatus
	}
	return res
}

// makeLikePostResponse LikePostメソッドのresponseを生成し返す
func (s server) makeLikePostResponse(statusCode string) *postservice.LikePostResponse {
	res := &postservice.LikePostResponse{}
	if statusCode != "" {
		responseStatus := &postservice.ResponseStatus{
			Code: statusCode,
		}
		res.Status = responseStatus
	}
	return res
}

// makeNotLikePostResponse NotLikePostメソッドのresponseを生成し返す
func (s server) makeNotLikePostResponse(statusCode string) *postservice.NotLikePostResponse {
	res := &postservice.NotLikePostResponse{}
	if statusCode != "" {
		responseStatus := &postservice.ResponseStatus{
			Code: statusCode,
		}
		res.Status = responseStatus
	}
	return res
}

// makeDeletePostResponse DeletePostメソッドのresponseを生成し返す
func (s server) makeDeletePostResponse(statusCode string) *postservice.DeletePostResponse {
	res := &postservice.DeletePostResponse{}
	if statusCode != "" {
		responseStatus := &postservice.ResponseStatus{
			Code: statusCode,
		}
		res.Status = responseStatus
	}
	return res
}
