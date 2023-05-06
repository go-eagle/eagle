# 🦅 eagle

 [![GitHub Workflow Status](https://img.shields.io/github/workflow/status/go-eagle/eagle/Go?style=flat-square)](https://github.com/go-eagle/eagle)
 [![codecov](https://codecov.io/gh/go-eagle/eagle/branch/master/graph/badge.svg)](https://codecov.io/gh/go-eagle/eagle)
 [![GolangCI](https://golangci.com/badges/github.com/golangci/golangci-lint.svg)](https://golangci.com)
 [![godoc](https://godoc.org/github.com/go-eagle/eagle?status.svg)](https://godoc.org/github.com/go-eagle/eagle)
 [![Gitter](https://badges.gitter.im/go-eagle/eagle.svg)](https://gitter.im/go-eagle/eagle?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
 <a href="http://opentracing.io"><img src="https://img.shields.io/badge/OpenTracing-enabled-blue.svg" alt="OpenTracing Badge"></a>
 [![Go Report Card](https://goreportcard.com/badge/github.com/go-eagle/eagle)](https://goreportcard.com/report/github.com/go-eagle/eagle)
 [![gitmoji](https://img.shields.io/badge/gitmoji-%20%F0%9F%98%9C%20%F0%9F%98%8D-FFDD67.svg?style=flat-square)](https://github.com/carloscuesta/gitmoji)
 [![License](https://img.shields.io/github/license/go-eagle/eagle?style=flat-square)](/LICENSE)

## Introduction

A Go framework suitable for rapid business development, which can quickly build API services or Web sites.

## Features

- Router [Gin](https://github.com/gin-gonic/gin) 
- Middleware [Gin](https://github.com/gin-gonic/gin) 
- Database [GORM](https://github.com/jinzhu/gorm)
- Document [Swagger](https://swagger.io/) 生成
- Config [Viper](https://github.com/spf13/viper)
- Auth [JWT](https://jwt.io/) 
- Validator [validator](https://github.com/go-playground/validator)
- Cron [cron](https://github.com/robfig/cron)
- Test [GoConvey](http://goconvey.co/)
- CI/CD [GitHub Actions](https://github.com/actions)
- Lint [GolangCI-lint](https://golangci.com/)

## Directory Structure

```shell
├── Makefile                     # 项目管理文件
├── api                          # grpc客户端和Swagger 文档
├── cmd                          # 脚手架目录
├── config                       # 配置文件统一存放目录
├── docs                         # 框架相关文档
├── internal                     # 业务目录
│   ├── cache                    # 基于业务封装的cache
│   ├── handler                  # http 接口
│   ├── middleware               # 自定义中间件
│   ├── model                    # 数据库 model
│   ├── dao                      # 数据访问层
│   ├── ecode                    # 业务自定义错误码
│   ├── routers                  # 业务路由
│   ├── server                   # http server 和 grpc server
│   └── service                  # 业务逻辑层
├── logs                         # 存放日志的目录
├── main.go                      # 项目入口文件
├── pkg                          # 公共的 package
├── test                         # 单元测试依赖的配置文件，主要是供docker使用的一些环境配置文件
└── scripts                      # 存放用于执行各种构建，安装，分析等操作的脚本
```

## Installtion CLI

```bash
# go >= 1.16
go install github.com/go-eagle/eagle/cmd/eagle@latest

# go < 1.16
go get github.com/go-eagle/eagle/cmd/eagle
```

## Quick Start

```bash
eagle new eagle-demo
# or 
eagle new github.com/foo/eagle-demo

# build
make build

# run
./eagle-demo
```

## Documentation

[https://go-eagle.org/](https://go-eagle.org/)

## CHANGELOG

- [CHANGELOG](https://github.com/go-eagle/eagle/blob/master/CHANGELOG.md)

## Who is using

- [1024casts](https://1024casts.com)
- [FastIM](https://github.com/1024casts/fastim)
- [Go-microservice](https://github.com/go-microservice)

## Discussion

- Issue: https://github.com/go-eagle/eagle/issues
- Gitter: https://gitter.im/go-eagle/eagle

## Stargazers over time

[![Stargazers over time](https://starchart.cc/go-eagle/eagle.svg)](https://starchart.cc/go-eagle/eagle)

## License

MIT. See the [LICENSE](LICENSE) file for details.
