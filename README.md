<h1 align="center">snake</h1>

<p align="center">
一款适合于快速开发业务的Go框架，主要是提供API服务。
</p> 

 [![GitHub Workflow Status](https://img.shields.io/github/workflow/status/1024casts/snake/Go?style=flat-square)](https://github.com/1024casts/snake)
 [![codecov](https://codecov.io/gh/1024casts/snake/branch/master/graph/badge.svg)](https://codecov.io/gh/1024casts/snake)
 [![GolangCI](https://golangci.com/badges/github.com/golangci/golangci-lint.svg)](https://golangci.com)
 [![godoc](https://godoc.org/github.com/1024casts/snake?status.svg)](https://godoc.org/github.com/1024casts/snake)
 [![Go Report Card](https://goreportcard.com/badge/github.com/1024casts/snake)](https://goreportcard.com/report/github.com/1024casts/snake)
 [![License](https://img.shields.io/github/license/1024casts/snake?style=flat-square)](/LICENSE)


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
- 支持多环境运行
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

// 运行
./snake -c conf/config.local.yaml
```

## 常用命令
 - make help 查看帮助
 - make dep 下载go依赖包
 - make build 编译项目
 - make swag-init 生成接口文档(需要重新编译)
 - make test-coverage 生成测试覆盖
 - make lint 检查代码规范

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
