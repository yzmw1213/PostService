package grpc

import (
	"context"
	"log"

	"github.com/yzmw1213/PostService/domain/model"
	"github.com/yzmw1213/PostService/grpc/postservice"
)

func (s server) CreatePost(ctx context.Context, req *postservice.CreatePostRequest) (*postservice.CreatePostResponse, error) {
	postData := req.GetPost()

	post := makePostModel(postData)

	post, err := s.PostUsecase.Create(post)
	if err != nil {
		return nil, err
	}
	res := &postservice.CreatePostResponse{
		Post: makeGrpcPost(post),
	}
	return res, nil

}

func (s server) DeletePost(ctx context.Context, req *postservice.DeletePostRequest) (*postservice.DeletePostResponse, error) {
	postData := req.GetId()
	post := &model.Post{
		ID: postData,
	}
	if err := s.PostUsecase.Delete(post); err != nil {
		return nil, err
	}
	res := &postservice.DeletePostResponse{}
	return res, nil
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
	row, err := s.PostUsecase.Read(ID)
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

	post := makePostModel(postData)

	updatedPost, err := s.PostUsecase.Update(post)
	if err != nil {
		return nil, err
	}
	res := &postservice.UpdatePostResponse{
		Post: makeGrpcPost(updatedPost),
	}
	return res, nil
}

func makePostModel(gUser *postservice.Post) *model.Post {
	post := &model.Post{
		ID:           gUser.GetId(),
		Title:        gUser.GetTitle(),
		Content:      gUser.GetContent(),
		CreateUserID: gUser.GetCreateUserId(),
	}
	return post
}

func makeGrpcPost(post *model.Post) *postservice.Post {
	gPost := &postservice.Post{
		Id:           post.ID,
		Title:        post.Title,
		Content:      post.Content,
		CreateUserId: post.CreateUserID,
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
