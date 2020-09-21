package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	pb "github.com/1024casts/snake/internal/rpc/user/v0"
)

func main() {
	serviceAddress := "127.0.0.1:1234"
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
	if err != nil {
		panic("connect error")
	}
	defer conn.Close()

	userClient := pb.NewUserServiceClient(conn)
	userReq := &pb.PhoneLoginRequest{
		Phone:      13010102020,
		VerifyCode: 123456,
	}
	reply, _ := userClient.LoginByPhone(context.Background(), userReq)
	fmt.Printf("UserService LoginByPhone : %+v", reply)
}
