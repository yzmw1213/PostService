package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/go-playground/validator/v10"

	"github.com/yzmw1213/GoMicroApp/domain/model"
	"github.com/yzmw1213/GoMicroApp/grpc/blog_grpc"
	"github.com/yzmw1213/GoMicroApp/usecase/interactor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	Usecase interactor.BlogInteractor
}

// NewBlogGrpcServer gRPCサーバー起動
func NewBlogGrpcServer() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	blogServer := &server{}

	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)

	blog_grpc.RegisterBlogServiceServer(s, blogServer)

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

func (s server) CreateBlog(ctx context.Context, req *blog_grpc.CreateBlogRequest) (*blog_grpc.CreateBlogResponse, error) {
	postData := req.GetBlog()
	validate := validator.New()

	blog := &model.Blog{
		AuthorID: postData.GetAuthorId(),
		Title:    postData.GetTitle(),
		Content:  postData.GetContent(),
	}

	// Blog構造体のバリデーション
	if error := validate.Struct(blog); error != nil {
		return nil, error
	}

	if err := s.Usecase.Create(blog); err != nil {
		return nil, err
	}
	res := &blog_grpc.CreateBlogResponse{
		Blog: postData,
	}
	return res, nil
}

func (s server) DeleteBlog(ctx context.Context, req *blog_grpc.DeleteBlogRequest) (*blog_grpc.DeleteBlogResponse, error) {
	postData := req.GetBlogId()
	blog := &model.Blog{
		BlogID: postData,
	}
	if err := s.Usecase.Delete(blog); err != nil {
		return nil, err
	}
	res := &blog_grpc.DeleteBlogResponse{}
	return res, nil
}

func (s server) ListBlog(req *blog_grpc.ListBlogRequest, stream blog_grpc.BlogService_ListBlogServer) error {
	rows, err := s.Usecase.List()
	if err != nil {
		return err
	}
	for _, blog := range rows {
		blog := &blog_grpc.Blog{
			BlogId:   blog.BlogID,
			AuthorId: blog.AuthorID,
			Title:    blog.Title,
			Content:  blog.Content,
		}
		res := &blog_grpc.ListBlogResponse{
			Blog: blog,
		}
		sendErr := stream.Send(res)
		if sendErr != nil {
			log.Fatalf("Error while sending response to client :%v", sendErr)
			return sendErr
		}
	}

	return nil
}

func (s server) ReadBlog(ctx context.Context, req *blog_grpc.ReadBlogRequest) (*blog_grpc.ReadBlogResponse, error) {
	blogID := req.GetBlogId()
	row, err := s.Usecase.Read(blogID)
	if err != nil {
		return nil, err
	}
	blog := &blog_grpc.Blog{
		BlogId:   row.BlogID,
		AuthorId: row.AuthorID,
		Title:    row.Title,
		Content:  row.Content,
	}
	res := &blog_grpc.ReadBlogResponse{
		Blog: blog,
	}
	return res, nil
}

func (s server) UpdateBlog(ctx context.Context, req *blog_grpc.UpdateBlogRequest) (*blog_grpc.UpdateBlogResponse, error) {
	postData := req.GetBlog()
	validate := validator.New()

	blog := &model.Blog{
		BlogID:   postData.GetBlogId(),
		AuthorID: postData.GetAuthorId(),
		Title:    postData.GetTitle(),
		Content:  postData.GetContent(),
	}

	// Blog構造体のバリデーション
	if error := validate.Struct(blog); error != nil {
		return nil, error
	}

	if err := s.Usecase.Update(blog); err != nil {
		return nil, err
	}
	res := &blog_grpc.UpdateBlogResponse{
		Blog: postData,
	}
	return res, nil
}
