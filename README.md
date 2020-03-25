<h1 align="center">snake</h1>

<p align="center">
一款适合于快速开发业务的Go框架，主要是提供API服务。
</p> 

<p align="center">
 <img alt="GitHub Workflow Status" src="https://img.shields.io/github/workflow/status/1024casts/snake/Go?style=flat-square">
 <a href="https://codecov.io/gh/1024casts/snake">
  <img src="https://codecov.io/gh/1024casts/snake/branch/master/graph/badge.svg" />
 </a>
 <a href="https://pkg.go.dev/github.com/1024casts/snake" rel="nofollow">
   <img src="https://camo.githubusercontent.com/cc8ece95e1e3139079889b8311ea73a4ab23b05f/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f676f2d72657374792f72657374793f7374617475732e737667" alt="GoDoc" data-canonical-src="https://godoc.org/github.com/1024casts/snake?status.svg" style="max-width:100%;">
 </a>
 <a href="https://goreportcard.com/report/github.com/1024casts/snake" rel="nofollow">
   <img src="https://camo.githubusercontent.com/4693f183dcbcd9ede6b5f53a7359f3e0709af676/68747470733a2f2f676f7265706f7274636172642e636f6d2f62616467652f6769746875622e636f6d2f737465696e666c6574636865722f61706974657374" alt="Go Report Card" data-canonical-src="https://goreportcard.com/badge/github.com/1024casts/snake" style="max-width:100%;">
 </a>
 <img alt="GitHub" src="https://img.shields.io/github/license/1024casts/snake?style=flat-square">
</p> 


## 技术栈

- 框架路由使用 [gin](https://github.com/gin-gonic/gin) 路由
- 中间件使用 [gin](https://github.com/gin-gonic/gin) 框架的中间件
- 数据库组件 [gorm](https://github.com/jinzhu/gorm)
- 文档使用 [swagger](https://swagger.io/) 生成
- 配置文件解析库 [viper](https://github.com/spf13/viper)
- 使用 [JWT](https://jwt.io/) 进行身份鉴权认证
- 校验器 [validator](https://gopkg.in/go-playground/validator.v9)  也是 gin 框架默认的校验器，当前最新是v9版本
- 任务调度 [cron](https://github.com/robfig/cron)
- 包管理工具 [go module](https://github.com/golang/go/wiki/Modules)
- 测试框架 [goConvey](http://goconvey.co/)
- CI/CD [Github Actions](https://github.com/actions)
- 容器化 [Docker](https://www.docker.com/)

## 特性

- 遵循 RESTful API 设计规范
- 基于 GIN WEB 框架，提供了丰富的中间件支持（用户认证、跨域、访问日志、请求频率限制、追踪 ID 等）
- 基于 GORM 的数据库存储
- JWT 认证
- 支持 Swagger 文档(基于[swaggo](https://github.com/swaggo/swag))
- 使用 make 来管理Go工程
- 使用 shell(admin.sh) 脚本来管理进程
- todo: 单元测试(基于net/http/httptest包，覆盖所有接口层的测试)
- todo: 使用github actions 进行CI/CD
- todo: 使用docker进行快速部署

## 目录结构

```shell
├── Makefile                     # 项目管理文件
├── admin.sh                     # 进程的start|stop|status|restart控制文件
├── conf                         # 配置文件统一存放目录
├── config                       # 专门用来处理配置和配置文件的Go package                 
├── db.sql                       # 在部署新环境时，可以登录MySQL客户端，执行source db.sql创建数据库和表
├── docs                         # swagger文档，执行 swag init 生成的
├── handler                      # 类似MVC架构中的C，用来读取输入，并将处理流程转发给实际的处理函数，最后返回结果
├── log                          # 存放日志的目录
├── main.go                      # 项目入口文件
├── model                        # 数据库model
├── pkg                          # 一些封装好的package
├── repository                   # 数据访问层
├── router                       # 路由及中间件目录
├── service                      # 业务逻辑封装
├── schedule                     # 任务调度配置目录
├── util                         # 业务工具包
└── wrktest.sh                   # API 性能测试脚本
```

## API风格和媒体类型

Go 语言中常用的 API 风格是 RPC 和 REST，常用的媒体类型是 JSON、XML 和 Protobuf。  
在 Go API 开发中常用的组合是 `gRPC + Protobuf` (更适合调用频繁的微服务场景) 和 `REST + JSON`。

本项目使用 API 风格采用 REST，媒体类型选择 JSON 。

<details>
 <summary>REST</summary>

REST 风格虽然适用于很多传输协议，但在实际开发中，REST 由于天生和 HTTP 协议相辅相成，因此 HTTP 协议已经成了实现 RESTful API 事实上的标准。  
在 HTTP 协议中通过 POST、DELETE、PUT、GET 方法来对应 REST 资源的增、删、改、查操作，具体的对应关系如下：

| HTTP方法 | 行为 | URI | 示例说明 |  
| :------ | :------ | :------ | :------ |
| GET	  | 获取资源列表  |	/users | 获取用户列表 |
| GET	  | 获取一个具体的资源  |	/users/admin |	获取 admin 用户的详细信息 |
| POST	  | 创建一个新的资源  | /users |	创建一个新用户 |
| PUT	  | 更新一个资源	| /users/1 |	更新 id 为 1 的用户 |
| DELETE  |	删除服务器上的一个资源	| /users/1 |	删除 id 为 1 的用户 |
</details>

## 错误码设计

[看这里](https://github.com/1024casts/snake/tree/master/pkg/errno)

## 启动项目

```bash
// 下载依赖
make dep

// 编译项目
make build

// 本地环境
cp config.sample.yaml config.local.yaml
./snake -c conf/config.local.yaml

// 线上环境类似操作
cp config.sample.yaml config.prod.yaml
./snake -c conf/config.prod.yaml
```

## 常用命令
 - make help 查看帮助
 - make dep 下载go依赖包
 - make build 编译项目
 - make swag-init 生成接口文档(需要重新编译)
 - make test-coverage 生成测试覆盖

## 模块
 - 用户(示例)
 
## 接口文档
`http://localhost:8080/swagger/index.html`

## 部署

### docker部署

Happy Coding. ^_^

## 谁在用
 - [1024课堂](https://1024casts.com)

## Discussion
- Issue: https://github.com/1024casts/snake/issues

## License
MIT. See the [LICENSE](LICENSE) file for details.
