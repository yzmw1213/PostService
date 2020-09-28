package grpc

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/yzmw1213/PostService/grpc/post_grpc"
	"github.com/yzmw1213/PostService/usecase/interactor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	PostUsecase interactor.PostInteractor
	TagUsecase  interactor.TagInteractor
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

func makeServer() *grpc.Server {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc.UnaryServerInterceptor(transmitStatusInterceptor)),
	)

	return s
}
