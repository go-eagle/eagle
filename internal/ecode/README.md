# 业务错误码定义

> 公共错误码已经在 `github.com/go-eagle/eagle/pkg/errno` 包中，可以直接使用

业务的错误码可以根据模块按文件进行定义

使用时公共错误码 以 `errno.开头`，业务错误码以 `ecode.开头`

## Demo

```go
// 公共错误码
import "github.com/go-eagle/eagle/pkg/errno"
...
errno.InternalServerError

// 业务错误码
import "github.com/go-eagle/eagle/internal/ecode"
...
ecode.ErrUserNotFound
```