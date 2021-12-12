# rpc server and client

## define a proto

```proto
syntax = "proto3";

package helloworld;

option go_package="github.com/go-eagle/eagle/examples/helloworld/helloworld";

// The greeting service definition.
service Greeter {
  // Sends a greeting
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}

// The request message containing the user's name.
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}
```

## generated helper code

the protocol buffer compiler generates codes that has

- message serialization code(*.pb.go)
- remote interface stub for Client to call with the methods(*_grpc.pb.go)
- abstract interface for Server code to implement(*_grpc.pb.go)

## try it out.

enter the project root directory

```bash
cd {root_path}

go get -v google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
go get -v google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
```

## gen proto

old-way
> 使用：github.com/golang/protobuf

```bash
$ protoc -I . --go_out=plugins=grpc,paths=source_relative:. examples/helloworld/protos/greeter.proto
```
> 生成的 `*.pb.go` 包含消息序列化代码和 `gRPC` 代码.

new-way

> 使用：google.golang.org/protobuf

```bash
$ protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    examples/helloworld/protos/greeter.proto
```
> 会生成两个文件 `*.pb.go` 和 `*._grpc.pb.go`, 分别是消息序列化代码和 `gRPC` 代码.

> 官方说明：https://github.com/protocolbuffers/protobuf-go/releases/tag/v1.20.0#v1.20-grpc-support

> https://grpc.io/docs/languages/go/quickstart/#regenerate-grpc-code

## Run

1. run the server

```bash
cd examples/helloworld/server
go run main.go
```

2. run the client from another terminal

```bash
cd examples/helloworld/client
go run main.go
```

3. You’ll see the following output:

```bash
Greeting : "Hello eagle"
```

## Reference

- https://grpc.io/docs/languages/go/quickstart/
- https://developers.google.com/protocol-buffers/docs/proto3
- https://grpc.io/docs/guides/error/
- https://github.com/avinassh/grpc-errors/blob/master/go/client.go
- https://stackoverflow.com/questions/64828054/differences-between-protoc-gen-go-and-protoc-gen-go-grpc
