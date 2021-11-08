package server

import (
	"log"

	"net"

	"google.golang.org/grpc"

	v1 "github.com/go-eagle/eagle/api/grpc/user/v1"
	"github.com/go-eagle/eagle/internal/service"
	"github.com/go-eagle/eagle/pkg/app"
)

// NewGRPCServer creates a gRPC server
func NewGRPCServer(cfg *app.ServerConfig) *grpc.Server {
	lis, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	v1.RegisterUserServiceServer(grpcServer, service.Svc.Users())
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve grpc server: %v", err)
	}
	log.Printf("serve grpc server is success, port:%s", cfg.Addr)

	return grpcServer
}
