# rpc server and client

## try it out.

enter the project root directory

```bash
cd {root_path}
```


### gen proto

```bash
protoc -I . --go_out=plugins=grpc,paths=source_relative:. examples/helloworld/protos/greeter.proto
```

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