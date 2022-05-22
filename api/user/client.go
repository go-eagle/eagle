package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "github.com/go-eagle/eagle/api/user/v1"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	conn, err := grpc.Dial(*addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	userClient := pb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	userReq := &pb.PhoneLoginRequest{
		Phone:      13010102020,
		VerifyCode: 123456,
	}
	reply, err := userClient.LoginByPhone(ctx, userReq)
	if err != nil {
		log.Fatalf("[rpc] user login by phone err: %v", err)
	}
	fmt.Printf("UserService LoginByPhone : %+v", reply)
}
