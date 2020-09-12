package main

import (
	"flag"
	"log"
	"net"

	"github.com/1024casts/snake/internal/service/user"

	pb "github.com/1024casts/snake/internal/rpc/user/v0"
	"google.golang.org/grpc"
)

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userService := user.NewUserService()
	pb.RegisterUserServiceServer(grpcServer, userService)
	grpcServer.Serve(lis)
}
