package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/yzmw1213/PostService/domain/model"
	"github.com/yzmw1213/PostService/grpc/post_grpc"
	"github.com/yzmw1213/PostService/usecase/interactor"

	"google.golang.org/grpc/reflection"
)

type server struct {
	Usecase interactor.PostInteractor
}

// NewPostGrpcServer gRPCサーバー起動
func NewPostGrpcServer() {
	lis, err := net.Listen("tcp", "0.0.0.0:50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	postServer := &server{}

	s := makeServer()

	post_grpc.RegisterPostServiceServer(s, postServer)

	// Register reflection service on gRPC server.
	reflection.Register(s)
	log.Println("main grpc server has started")

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a sgnal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the client")
	lis.Close()
	fmt.Println("End of Program")

}

func (s server) CreatePost(ctx context.Context, req *post_grpc.CreatePostRequest) (*post_grpc.CreatePostResponse, error) {
	postData := req.GetPost()

	post := makeModel(postData)

	post, err := s.Usecase.Create(post)
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
	if err := s.Usecase.Delete(post); err != nil {
		return nil, err
	}
	res := &post_grpc.DeletePostResponse{}
	return res, nil
}

func (s server) ListPost(req *post_grpc.ListPostRequest, stream post_grpc.PostService_ListPostServer) error {
	rows, err := s.Usecase.List()
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
	row, err := s.Usecase.Read(ID)
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

	post := makeModel(postData)

	updatedPost, err := s.Usecase.Update(post)
	if err != nil {
		return nil, err
	}
	res := &post_grpc.UpdatePostResponse{
		Post: makeGrpcPost(updatedPost),
	}
	return res, nil
}

func makeModel(gUser *post_grpc.Post) *model.Post {
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
