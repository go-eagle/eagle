# server

接口实现层，可以理解成 MVC 的控制器层。

## 目录结构

server 层跟 rpc 目录一一对应，原则上使用 `{service_name}server{version_number}` 方式命名。
例如，`rpc/user/v0` 服务对应的 server 目录是 `server/userserver0`。

一个目录下可以有多个 *.go 文件，每个 service 一个。

## 实现接口

服务接口定义在 rpc 目录对应的 echo.twirp.go 中，是自动生成的。

```go
package userserver0

import (
	// 标准库单列一组
	"context"

	// 框架库单列一组
	"github.com/1024casts/snake/internal/repository/login"
	"github.com/1024casts/snake/pkg/util/conf"

	// pb 定义单列一组
	pb "github.com/1024casts/snake/internal/rpc/user/v0"
)

// 服务对象，约定为 Server
type EchoServer struct{}

// 接口实现，三步走：处理入参、调用服务、返回出参
func (s *EchoServer) ClearLoginCache(ctx context.Context, req *pb.ClearRequest) (*pb.EmptyReply, error) {
	// 处理入参
	mid := req.GetMid()

	// 调用 service 层或者 dao 层完成业务逻辑
	login.ClearUID(ctx, mid)

	// 返回出参
	reply := &pb.EmptyReply{}

	return reply, nil
}
```

## 注册服务

请参考 [cmd/server/README.md](../cmd/server/README.md)。

## 错误处理

### 异常/错误

**错误** 是 __计划内__ 的情形，例如用户输入密码不匹配、用户余额不足等等。
**异常** 是 __计划外__ 的情形，例如用户提交的参数类型跟接口定义不匹配、DB 连接超时等等。

**错误** 可以认为是一种特殊的“正常情况”, **异常** 则是真正的“不正常情况”。

### 处理错误

客户端需要根据不同业务需求处理 **错误**, 例如用户未登录则需要跳转到登录页面。所以，我需要使用错误码来返回错误信息。

处理代码示例如下：

```go
resp := &pb.Resp{}

resp.Code = 100
resp.Msg = "Need Login"

return nil, resp
```

以上代码会返回如下 HTTP 信息：
```
HTTP/1.1 200 OK
Content-Length: 355
Content-Type: application/json
Date: Tue, 14 Aug 2018 03:05:41 GMT
X-Trace-Id: 3kclnknyzmamo

{
    "code": 100,
    "msg": "Need Login",
    "data": {}
}
```

### 处理异常

正常的客户端会严格按照接口定义调用接口，只有客户端有 bug 或者服务端有问题的时候才会遇到 **异常**。
在这种情况下，首先，我们无法从错误中恢复；其次，这类错误的处理方式跟具体的业务没有关系的；最后，我们需要 **及时发现** 这类问题并修复。
所以，我们需要使用 HTTP 的 4xx 和 5xx 状态码来返回错误信息。

处理代码示例如下：

```go
import "github.com/1024casts/snake/pkg/errno"
// ...

// 这是客户端问题，返回 HTTP 4xx 状态码
if req.ID <= 0 {
	return nil, errors.InvalidArgumentError("id", "must > 0")
}

// HTTP/1.1 400 Bad Request
// Content-Length: 104
// Content-Type: application/json
// Date: Tue, 14 Aug 2018 03:09:30 GMT
// X-Trace-Id: kg1od386gjto
//
// {
//     "code": "invalid_argument",
//     "meta": {
//         "argument": "page_size"
//     },
//     "msg": "page_size page_size must be > 0"
// }

// 这是服务端问题，返回 HTTP 5xx 状态码
if err := bookshelf.AddFavorite(ctx, id); err != nil {
	return nil, err
}

// HTTP/1.1 500 Internal Server Error
// Content-Length: 112
// Content-Type: application/json
// Date: Wed, 15 Aug 2018 08:50:47 GMT
// X-Trace-Id: 3njq5120j3c1n
//
// {
//     "code": "internal",
//     "meta": {
//         "cause": "*net.OpError"
//     },
//     "msg": "dial tcp :0: connect: can't assign requested address"
// }
```

我们可以通过 SLB 报警及时发现此类错误并减少业务损失。