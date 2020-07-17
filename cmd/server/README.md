
## 注册服务

- 自动注册服务请参考 [rpc/README.md](https://github.com/1024casts/snake/blob/master/internal/rpc/README.md)  
- 实现服务接口请参考 [server/README.md](https://github.com/1024casts/snake/blob/master/internal/server/README.md)。

## 启动服务

```bash
# 对外服务
go run main.go server --port=8080
# 对内服务
go run main.go server --port=8080 --internal
```