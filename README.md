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

A Go framework suitable for rapid business development, which can quickly build API services or Web sites.   
English | [中文文档](https://github.com/go-eagle/eagle/blob/master/README_ZH.md)

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
