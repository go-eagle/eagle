package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/go-eagle/eagle/examples/helloworld/protos"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	protos.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *protos.HelloRequest) (*protos.HelloReply, error) {
	if in.Name == "" {
		return nil, fmt.Errorf("invalid argument %s", in.Name)
	}
	return &protos.HelloReply{Message: fmt.Sprintf("Hello %+v", in.Name)}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	protos.RegisterGreeterServer(grpcServer, &server{})
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve grpc server: %v", err)
	}
	log.Printf("serve grpc server is success, port:%s", "9090")
}
