# 🦅 eagle

 [![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/go-eagle/eagle/test.yml?branch=master&style=flat-square)](https://github.com/go-eagle/eagle)
 [![codecov](https://codecov.io/gh/go-eagle/eagle/branch/master/graph/badge.svg)](https://codecov.io/gh/go-eagle/eagle)
 [![GolangCI](https://golangci.com/badges/github.com/golangci/golangci-lint.svg)](https://golangci.com)
 [![godoc](https://godoc.org/github.com/go-eagle/eagle?status.svg)](https://godoc.org/github.com/go-eagle/eagle)
 [![Gitter](https://badges.gitter.im/go-eagle/eagle.svg)](https://gitter.im/go-eagle/eagle?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
 <a href="http://opentracing.io"><img src="https://img.shields.io/badge/OpenTracing-enabled-blue.svg" alt="OpenTracing Badge"></a>
 [![Go Report Card](https://goreportcard.com/badge/github.com/go-eagle/eagle)](https://goreportcard.com/report/github.com/go-eagle/eagle)
 [![gitmoji](https://img.shields.io/badge/gitmoji-%20%F0%9F%98%9C%20%F0%9F%98%8D-FFDD67.svg?style=flat-square)](https://github.com/carloscuesta/gitmoji)
 [![License](https://img.shields.io/github/license/go-eagle/eagle?style=flat-square)](/LICENSE)

A Go framework suitable for rapid business development, which can quickly build API services or Web sites.   
English | [中文文档](https://github.com/go-eagle/eagle/blob/master/README_ZH.md)

## Features

- API framework [gin](https://github.com/gin-gonic/gin) 
- RPC framework [gRPC](https://github.com/grpc/grpc-go)
- Configuration [viper](https://github.com/spf13/viper)
- Logging component [zap](https://github.com/uber-go/zap)
- Database ORM component [gorm](https://github.com/go-gorm/gorm) | [MongoDB](https://github.com/mongodb/mongo-go-driver)
- Search component [Elasticsearch](https://github.com/elastic/go-elasticsearch)
- Cache component [go-redis](https://github.com/go-redis/redis), [ristretto](https://github.com/dgraph-io/ristretto)
- Message Queue [Rabbitmq](https://github.com/rabbitmq/amqp091-go) | [redis](https://github.com/hibiken/asynq)
- Authentication [JWT](https://jwt.io/) 
- Parameter Validator [validator](https://github.com/go-playground/validator)
- Scheduled tasks [cron](https://github.com/robfig/cron)
- Metrics monitoring [prometheus](https://github.com/prometheus/client_golang/prometheus), [grafana](https://github.com/grafana/grafana)
- Distributed Tracing [opentelemetry](https://github.com/open-telemetry/opentelemetry-go)
- Service registration and discovery [etcd](https://github.com/etcd-io/etcd) | [consul](https://github.com/hashicorp/consul) | [nacos](https://github.com/alibaba/nacos)
- Unit Test [GoConvey](http://goconvey.co/)
- Lint [GolangCI-lint](https://golangci.com/)
- CI/CD [GitHub Actions](https://github.com/actions), [docker](https://www.docker.com/), [kubernetes](https://github.com/kubernetes/kubernetes)

## Framework Layered Architecture
![eagle-framework-diagram](https://github.com/go-eagle/eagle/assets/3043638/cd05f6d5-058c-4ab0-87ee-47148e0c68aa)

## Logic Layered Architecture

Eagle utilizes a classic layered structure and employs the Wire dependency injection framework to enhance modularity and reduce coupling between components.

[![Leagle Layout Arch](https://raw.githubusercontent.com/go-eagle/eagle/master/docs/images/eagle-layout-arch.png)](https://starchart.cc/go-eagle/eagle)

## Directory Structure

```shell
├── Makefile                     
├── api                          
├── cmd                          
├── config                       
├── docs                         
├── internal                     
│   ├── cache                    
│   ├── handler                  
│   ├── middleware               
│   ├── model                    
│   ├── dao                      
│   ├── ecode                    
│   ├── routers                  
│   ├── server                   
│   └── service                  
├── logs                         
├── main.go                      
├── pkg                          
├── test                         
└── scripts                      
```

## Installtion CLI

```bash
GOPROXY="https://goproxy.cn,direct"

# go >= 1.16
go install github.com/go-eagle/eagle/cmd/eagle@latest

# go < 1.16
go get github.com/go-eagle/eagle/cmd/eagle
```

## Quick Start

```bash
# gen a server with http and gRPC
eagle new eagle-demo
# or 
eagle new github.com/foo/eagle-demo

# install dependence
go mod tidy

# run
make run
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
- Discord: https://discord.com/channels/968369660900814869

## Microservice Roadmap

![Microservice-roadmap](https://github.com/go-eagle/eagle/assets/3043638/c7ef237e-e0f9-4699-843d-54588b2bcec8)

## Contributing

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement". Don't forget to give the project a star! Thanks again!

- Fork the repository to your own GitHub account.
- Create a new branch for your changes.
- Make your changes to the code.
- Commit your changes and push the branch to your forked repository.
- Open a pull request on our repository.

## Stargazers over time

[![Stargazers over time](https://starchart.cc/go-eagle/eagle.svg)](https://starchart.cc/go-eagle/eagle)

## License

MIT. See the [LICENSE](LICENSE) file for details.
