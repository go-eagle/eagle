package main

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	pb "github.com/1024casts/snake/internal/rpc/user/v0"
	"github.com/1024casts/snake/pkg/log"
)

func main() {
	// todo: 配置可以从环境变量或配置文件读取
	serviceAddress := "127.0.0.1:1234"
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Sprintf("grpc dial err: %v", err)
		panic("grpc dial err")
	}
	defer conn.Close()

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
