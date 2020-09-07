package grpc

import "google.golang.org/grpc"

func makeServer() *grpc.Server {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc.UnaryServerInterceptor(transmitStatusInterceptor)),
	)

	return s
}
