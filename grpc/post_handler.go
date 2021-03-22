package grpc

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/yzmw1213/PostService/aws"
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
	// StatusCreateCommentSuccess コメント作成成功ステータス
	StatusCreateCommentSuccess string = "COMMENT_CREATE_SUCCESS"
	// StatusUpdateCommentSuccess コメント更新成功ステータス
	StatusUpdateCommentSuccess string = "COMMENT_UPDATE_SUCCESS"
	// StatusDeletePostsCommentsByUserIDSuccess 投稿ユーザーID指定削除成功ステータス
	StatusDeletePostsCommentsByUserIDSuccess string = "COMMENT_DELETEE_BY_USERID_SUCCESS"
	// StatusDeletePostSuccess 投稿削除成功ステータス
	StatusDeletePostSuccess string = "POST_DELETE_SUCCESS"
	// StatusPostNotExists 指定した投稿の登録がない時のエラーステータス
	StatusPostNotExists string = "POST_NOT_EXISTS_ERROR"
	// StatusPostTitleStringCount 件名文字数が無効のエラーステータス
	StatusPostTitleStringCount string = "POST_TITLE_COUNT_ERROR"
	// StatusPostContentStringCount 投稿内容文字数が無効のエラーステータス
	StatusPostContentStringCount string = "POST_CONTENT_COUNT_ERROR"
	// StatusCommentContentStringCount 投稿内容文字数が無効のエラーステータス
	StatusCommentContentStringCount string = "POST_CONTENT_COUNT_ERROR"
)

var (
	bucket = os.Getenv("AWS_S3_BUCKET_NAME")
	region = os.Getenv("AWS_S3_REGION")
)

func (s server) CreatePost(ctx context.Context, req *postservice.CreatePostRequest) (*postservice.CreatePostResponse, error) {
	var location string
	postData := req.GetPost()

	post := makePostModel(postData)
	tags := makePostTagModel(postData)

	joinPost := &model.JoinPost{
		Post:     post,
		PostTags: tags,
	}

	if isBase64(post.Image) == true {
		// 画像をS3にアップロードし、URLを受け取る。
		location, err = aws.Upload(post.Image[strings.IndexByte(post.Image, ',')+1:])

		if err != nil {
			return nil, err
		}
		joinPost.Post.Image = location
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
	var posts []*postservice.Post
	condition := req.GetCondition()
	ID := req.GetId()

	rows, err := s.PostUsecase.List(condition, ID)
	if err != nil {
		return s.makeListPostResponse(posts), err
	}
	for _, post := range rows {
		post := makeGrpcPost(&post)
		posts = append(posts, post)
	}
	return s.makeListPostResponse(posts), nil
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
	// 更新時はimageの更新は行わない
	joinPost.Post.Image = ""

	if _, err := s.PostUsecase.Update(joinPost); err != nil {
		return nil, err
	}

	return s.makeUpdatePostResponse(StatusUpdatePostSuccess), nil
}

func (s server) LikePost(ctx context.Context, req *postservice.LikePostRequest) (*postservice.LikePostResponse, error) {
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
	postLikeUser := &model.PostLikeUser{
		PostID: req.GetId(),
		UserID: req.GetUserId(),
	}

	if _, err := s.PostUsecase.NotLike(postLikeUser); err != nil {
		return nil, err
	}
	return s.makeNotLikePostResponse(StatusNotLikePostSuccess), nil
}

func (s server) CreateComment(ctx context.Context, req *postservice.CreateCommentRequest) (*postservice.CreateCommentResponse, error) {
	comment := makeComment(req.Comment)
	if _, err := s.PostUsecase.CreateComment(comment); err != nil {
		return nil, err
	}
	return s.makeCreateCommentResponse(StatusCreateCommentSuccess), nil
}

func (s server) UpdateComment(ctx context.Context, req *postservice.UpdateCommentRequest) (*postservice.UpdateCommentResponse, error) {
	comment := makeComment(req.Comment)
	if _, err := s.PostUsecase.UpdateComment(comment); err != nil {
		return nil, err
	}
	return s.makeUpdateCommentResponse(StatusUpdateCommentSuccess), nil
}

func (s server) DeleteComment(ctx context.Context, req *postservice.DeleteCommentRequest) (*postservice.DeleteCommentResponse, error) {
	id := req.GetId()

	if err := s.PostUsecase.DeleteComment(id); err != nil {
		return nil, err
	}
	return s.makeDeleteCommentResponse(StatusDeletePostSuccess), nil
}

func (s server) DeletePostsCommentsByUserID(ctx context.Context, req *postservice.DeletePostsCommentsByUserIDRequest) (*postservice.DeletePostsCommentsByUserIDResponse, error) {
	createUserID := req.GetCreateUserId()

	// 退会ユーザーの投稿記事を削除
	if err := s.PostUsecase.DeletePostsByUserID(createUserID); err != nil {
		return nil, err
	}

	// 退会ユーザーのコメントを削除
	if err := s.PostUsecase.DeleteCommentsByUserID(createUserID); err != nil {
		return nil, err
	}

	res := &postservice.DeletePostsCommentsByUserIDResponse{
		Status: &postservice.ResponseStatus{
			Code: StatusDeletePostsCommentsByUserIDSuccess,
		},
	}

	return res, nil
}

func makePostModel(gPost *postservice.Post) *model.Post {
	post := &model.Post{
		ID: gPost.GetId(),
		// Status:       gPost.GetStatus(),
		Title:        gPost.GetTitle(),
		Content:      gPost.GetContent(),
		Image:        gPost.GetImage(),
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
	return postTags
}

func makeComment(gComment *postservice.Comment) *model.Comment {
	comment := &model.Comment{
		CommentID:      gComment.Id,
		PostID:         gComment.PostId,
		CreateUserID:   gComment.CreateUserId,
		CommentContent: gComment.Content,
	}

	return comment
}

func makeGrpcPost(post *model.JoinPost) *postservice.Post {
	var tags []uint32
	var likeUsers []uint32
	var postComments []*postservice.Comment
	bucketEndpoint := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/", bucket, region)
	gPost := &postservice.Post{
		Id: post.Post.ID,
		// Status:       post.Status,
		Title:          post.Post.Title,
		Content:        post.Post.Content,
		Image:          fmt.Sprintf("%s%s", bucketEndpoint, post.Post.Image),
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

	// コメント
	for _, comment := range post.Comments {
		gcomment := makeGrpcComment(comment)
		postComments = append(postComments, gcomment)
	}
	gPost.Comments = postComments

	return gPost
}

func makeGrpcComment(jc model.JoinComment) *postservice.Comment {
	return &postservice.Comment{
		Id:             jc.Comment.CommentID,
		PostId:         jc.Comment.PostID,
		CreateUserId:   jc.Comment.CreateUserID,
		CreateUserName: jc.CreateUser.UserName,
		Content:        jc.Comment.CommentContent,
	}
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

// makeListPostResponse ListPostメソッドのresponseを生成し返す
func (s server) makeListPostResponse(posts []*postservice.Post) *postservice.ListPostResponse {
	res := &postservice.ListPostResponse{
		Count: uint32(len(posts)),
		Post:  posts,
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

// makeCreateCommentResponse CreateCommentメソッドのresponseを生成し返す
func (s server) makeCreateCommentResponse(statusCode string) *postservice.CreateCommentResponse {
	res := &postservice.CreateCommentResponse{}
	if statusCode != "" {
		responseStatus := &postservice.ResponseStatus{
			Code: statusCode,
		}
		res.Status = responseStatus
	}
	return res
}

// makeUpdateCommentResponse UpdateCommentメソッドのresponseを生成し返す
func (s server) makeUpdateCommentResponse(statusCode string) *postservice.UpdateCommentResponse {
	res := &postservice.UpdateCommentResponse{}
	if statusCode != "" {
		responseStatus := &postservice.ResponseStatus{
			Code: statusCode,
		}
		res.Status = responseStatus
	}
	return res
}

// makeDeleteCommentResponse DeleteCommentメソッドのresponseを生成し返す
func (s server) makeDeleteCommentResponse(statusCode string) *postservice.DeleteCommentResponse {
	res := &postservice.DeleteCommentResponse{}
	if statusCode != "" {
		responseStatus := &postservice.ResponseStatus{
			Code: statusCode,
		}
		res.Status = responseStatus
	}
	return res
}

// isBase64 base64データであるか判定
func isBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s[strings.IndexByte(s, ',')+1:])
	return err == nil
}
