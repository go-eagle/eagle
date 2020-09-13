package server

import (
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/1024casts/snake/internal/rpc/user/v0"
	"github.com/1024casts/snake/internal/service"
)

func New(svc *service.Service) {
	flag.Parse()

	// todo: get addr from conf
	lis, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, svc.UserSvc())
	grpcServer.Serve(lis)
}
