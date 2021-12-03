# rpc server and client

## try it out.

enter the project root directory

```bash
cd {root_path}
go get -v google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
go get -v google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
```

### gen proto

```bash
// old
$ protoc -I . --go_out=plugins=grpc,paths=source_relative:. examples/helloworld/protos/greeter.proto

// new
$ protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    examples/helloworld/protos/greeter.proto
```

> https://grpc.io/docs/languages/go/quickstart/#regenerate-grpc-code

### run server

```bash
cd examples/helloworld/server
go run main.go
```

### run client

```bash
cd examples/helloworld/client
go run main.go
```

### run result from client

```bash
service SayHello : message:"Hello eagle"
```

## Reference

- https://grpc.io/docs/languages/go/quickstart/
- https://developers.google.com/protocol-buffers/docs/proto3