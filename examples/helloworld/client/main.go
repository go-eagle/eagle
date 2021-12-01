package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-eagle/eagle/examples/helloworld/protos"

	"github.com/go-eagle/eagle/pkg/log"
	"google.golang.org/grpc"
)

func main() {
	serviceAddress := "127.0.0.1:1234"
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Printf("grpc dial err: %v", err)
		panic("grpc dial err")
	}
	defer func() {
		_ = conn.Close()
	}()

	cli := protos.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &protos.HelloRequest{
		Name: "eagle",
	}
	reply, err := cli.SayHello(ctx, req)
	if err != nil {
		log.Errorf("[rpc] user login by phone err: %v", err)
	}
	fmt.Printf("UserService LoginByPhone : %+v", reply)
}
