package main

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	"github.com/go-eagle/eagle/api/grpc/user/v1"
	"github.com/go-eagle/eagle/pkg/log"
)

func main() {
	// todo: 配置可以从环境变量或配置文件读取
	serviceAddress := "127.0.0.1:1234"
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Printf("grpc dial err: %v", err)
		panic("grpc dial err")
	}
	defer func() {
		_ = conn.Close()
	}()

	userClient := v1.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	userReq := &v1.PhoneLoginRequest{
		Phone:      13010102020,
		VerifyCode: 123456,
	}
	reply, err := userClient.LoginByPhone(ctx, userReq)
	if err != nil {
		log.Errorf("[rpc] user login by phone err: %v", err)
	}
	fmt.Printf("UserService LoginByPhone : %+v", reply)
}
