package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/jinzhu/gorm"
	"github.com/yzmw1213/GoMicroApp/db"
	"github.com/yzmw1213/GoMicroApp/domain/model"
	"github.com/yzmw1213/GoMicroApp/grpc/blog_grpc"
	"github.com/yzmw1213/GoMicroApp/usecase/interactor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	Usecase interactor.BlogInteractor
}

func NewBlogGrpcServer() {
	fmt.Println("Hello")
	lis, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	blogServer := &server{}

	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)

	blog_grpc.RegisterBlogServiceServer(s, blogServer)

	// Register reflection service on gRPC server.
	reflection.Register(s)

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
	blog := &model.Blog{
		AuthorId: postData.GetAuthorId(),
		Title:    postData.GetTitle(),
		Content:  postData.GetContent(),
	}
	if err := s.Usecase.CreateBlog(blog); err != nil {
		return nil, err
	}
	res := &blog_grpc.CreateBlogResponse{
		Blog: postData,
	}
	return res, nil
}

func (s server) ListBlog(req *blog_grpc.ListBlogRequest, stream blog_grpc.BlogService_ListBlogServer) error {
	rows, err := s.Usecase.ListBlog()
	if err != nil {
		return err
	}
	for _, blog := range rows {
		blog := &blog_grpc.Blog{
			AuthorId: blog.AuthorId,
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
func GetDB() *gorm.DB {
	return db.GetDB()
}
