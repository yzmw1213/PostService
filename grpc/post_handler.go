package grpc

import (
	"context"
	"log"

	"github.com/yzmw1213/PostService/domain/model"
	"github.com/yzmw1213/PostService/grpc/post_grpc"
)

func (s server) CreatePost(ctx context.Context, req *post_grpc.CreatePostRequest) (*post_grpc.CreatePostResponse, error) {
	postData := req.GetPost()

	post := makePostModel(postData)

	post, err := s.PostUsecase.Create(post)
	if err != nil {
		return nil, err
	}
	res := &post_grpc.CreatePostResponse{
		Post: makeGrpcPost(post),
	}
	return res, nil

}

func (s server) DeletePost(ctx context.Context, req *post_grpc.DeletePostRequest) (*post_grpc.DeletePostResponse, error) {
	postData := req.GetId()
	post := &model.Post{
		ID: postData,
	}
	if err := s.PostUsecase.Delete(post); err != nil {
		return nil, err
	}
	res := &post_grpc.DeletePostResponse{}
	return res, nil
}

func (s server) ListPost(req *post_grpc.ListPostRequest, stream post_grpc.PostService_ListPostServer) error {
	rows, err := s.PostUsecase.List()
	if err != nil {
		return err
	}
	for _, post := range rows {
		post := &post_grpc.Post{
			Id: post.ID,
		}
		res := &post_grpc.ListPostResponse{
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

func (s server) ReadPost(ctx context.Context, req *post_grpc.ReadPostRequest) (*post_grpc.ReadPostResponse, error) {
	ID := req.GetId()
	row, err := s.PostUsecase.Read(ID)
	if err != nil {
		return nil, err
	}
	post := &post_grpc.Post{
		Id:      row.ID,
		UserId:  row.UserID,
		Content: row.Content,
	}
	res := &post_grpc.ReadPostResponse{
		Post: post,
	}
	return res, nil
}

func (s server) UpdatePost(ctx context.Context, req *post_grpc.UpdatePostRequest) (*post_grpc.UpdatePostResponse, error) {
	postData := req.GetPost()

	post := makePostModel(postData)

	updatedPost, err := s.PostUsecase.Update(post)
	if err != nil {
		return nil, err
	}
	res := &post_grpc.UpdatePostResponse{
		Post: makeGrpcPost(updatedPost),
	}
	return res, nil
}

func makePostModel(gUser *post_grpc.Post) *model.Post {
	post := &model.Post{
		ID:      gUser.GetId(),
		UserID:  gUser.GetUserId(),
		Content: gUser.GetContent(),
	}
	return post
}

func makeGrpcPost(post *model.Post) *post_grpc.Post {
	gPost := &post_grpc.Post{
		Id:      post.ID,
		UserId:  post.UserID,
		Content: post.Content,
	}
	return gPost
}

func createPostRequest(post *post_grpc.Post) *post_grpc.CreatePostRequest {
	return &post_grpc.CreatePostRequest{
		Post: post,
	}
}

func updatePostRequest(post *post_grpc.Post) *post_grpc.UpdatePostRequest {
	return &post_grpc.UpdatePostRequest{
		Post: post,
	}
}
