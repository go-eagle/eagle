package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"google.golang.org/grpc"

	pb "github.com/go-eagle/eagle/examples/helloworld/helloworld"
	"github.com/go-eagle/eagle/pkg/log"
)

const (
	defaultName = "eagle"
)

var (
	addr = flag.String("addr", "localhost:9090", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	conn, err := grpc.Dial(*addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Printf("grpc dial err: %v", err)
		panic("grpc dial err")
	}
	defer func() {
		_ = conn.Close()
	}()

	cli := pb.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.HelloRequest{
		Name: *name,
	}
	reply, err := cli.SayHello(ctx, req)
	if err != nil {
		log.Errorf("[rpc] SayHello err: %v", err)
	}
	fmt.Printf("service SayHello : %+v", reply)
}
