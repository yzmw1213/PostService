package main

import (
	"github.com/yzmw1213/PostService/db"
	"github.com/yzmw1213/PostService/grpc"
)

func main() {
	start()
}

func start() {
	db.Init()
	grpc.NewPostGrpcServer()
	defer db.Close()
}
