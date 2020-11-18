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

func (s server) ListPost(req *postservice.ListPostRequest, stream postservice.PostService_ListPostServer) error {
	rows, err := s.PostUsecase.List()
	if err != nil {
		return err
	}
	for _, post := range rows {
		post := &postservice.Post{
			Id: post.ID,
		}
		res := &postservice.ListPostResponse{
			Post: post,
		}
		sendErr := stream.Send(res)
		if sendErr != nil {
			log.Fatalf("Error while sending response to client :%v", sendErr)
			return sendErr
		}
	}

	return nil
}

func (s server) ReadPost(ctx context.Context, req *postservice.ReadPostRequest) (*postservice.ReadPostResponse, error) {
	ID := req.GetId()
	row, err := s.PostUsecase.GetByID(ID)
	if err != nil {
		return nil, err
	}
	post := &postservice.Post{
		Id:           row.ID,
		CreateUserId: row.CreateUserID,
		Content:      row.Content,
	}
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
	if _, err := s.PostUsecase.Update(joinPost.Post); err != nil {
		return nil, err
	}

	return s.makeUpdatePostResponse(StatusUpdatePostSuccess), nil
}

func makePostModel(gPost *postservice.Post) *model.Post {
	post := &model.Post{
		ID: gPost.GetId(),
		// Status:       gPost.GetStatus(),
		Title:        gPost.GetTitle(),
		Content:      gPost.GetContent(),
		MaxNum:       gPost.GetMaxNum(),
		Gender:       gPost.GetGender(),
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

func makeGrpcPost(post *model.Post) *postservice.Post {
	gPost := &postservice.Post{
		Id: post.ID,
		// Status:       post.Status,
		Title:        post.Title,
		Content:      post.Content,
		MaxNum:       post.MaxNum,
		Gender:       post.Gender,
		CreateUserId: post.CreateUserID,
		UpdateUserId: post.UpdateUserID,
	}
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