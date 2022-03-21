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
  // 字段定义，如果使用 form 表单传输，则只支持
  // int32, int64, uint32, unint64, double, float, bool, string,
  // 以及对应的 repeated 类型，不支持 map 和 message 类型！
  // 框架会自动解析并转换参数类型
  // 如果用 json 或 protobuf 传输则没有限制
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}

// 出参定义, 理论上可以输出任意消息
// 但我们约定只能包含 code, msg, data 三个字段，
// 其中 data 需要定义成 message
// 开源版本可以怱略这一约定
message HelloResponse {
  // 业务错误码[机读]，必须大于零
  // 小于零的主站框架在用，注意避让。
  int32 code = 1;
  // 业务错误信息[人读]
  string msg = 2;
  // 业务数据对象
  HelloReply data = 3;
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

`github.com/golang/protobuf` 和 `google.golang.org/protobuf` 有什么区别？https://developers.google.com/protocol-buffers/docs/reference/go/faq

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
- https://jbrandhorst.com/post/grpc-errors/
- https://godoc.org/google.golang.org/genproto/googleapis/rpc/errdetails
- https://cloud.google.com/apis/design/errors
- https://github.com/grpc/grpc/blob/master/doc/health-checking.md
- https://eddycjy.com/posts/where-is-proto/
- https://stackoverflow.com/questions/52969205/how-to-assert-grpc-error-codes-client-side-in-go
