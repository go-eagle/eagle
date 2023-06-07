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

一款适合于快速开发业务的 Go 框架，可快速构建 API 服务 或 Web 网站。

[English](https://github.com/go-eagle/eagle/blob/master/README.md) | 中文文档

## 官方文档

 - 开发文档 [https://go-eagle.org/](https://go-eagle.org/)

**Pro Tip:** 每个目录下基本都有 `README`，可以让框架使用起来更轻松 ^_^

## 设计思想和原则

框架中用到的设计思想和原则，尽量满足 "高内聚、低耦合"，主要遵从下面几个原则
- 1. 单一职责原则
- 2. 基于接口而非实现编程
- 3. 依赖注入
- 4. 多用组合
- 5. 迪米特法则

> 迪米特法则: 不该有直接依赖关系的类之间，不要有依赖；有依赖关系的类之间，尽量只依赖必要的接口

## ✨ 技术栈

- 框架路由使用 [Gin](https://github.com/gin-gonic/gin) 路由
- 中间件使用 [Gin](https://github.com/gin-gonic/gin) 框架的中间件
- 数据库组件 [GORM](https://github.com/jinzhu/gorm)
- 文档使用 [Swagger](https://swagger.io/) 生成
- 配置文件解析库 [Viper](https://github.com/spf13/viper)
- 使用 [JWT](https://jwt.io/) 进行身份鉴权认证
- 校验器使用 [validator](https://github.com/go-playground/validator)  也是 Gin 框架默认的校验器
- 任务调度 [cron](https://github.com/robfig/cron)
- 包管理工具 [Go Modules](https://github.com/golang/go/wiki/Modules)
- 测试框架 [GoConvey](https://github.com/smartystreets/goconvey)
- CI/CD [GitHub Actions](https://github.com/actions)
- 使用 [GolangCI-lint](https://golangci.com/) 进行代码检测
- 使用 make 来管理 Go 工程
- 使用 shell(admin.sh) 脚本来管理进程
- 使用 YAML 文件进行多环境配置

## 📗 目录结构

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

## 🛠️ 快速开始

### 方式一

直接Clone项目的方式，文件比较全

TIPS: 需要本地安装MySQL数据库和 Redis

```bash
# 下载安装，可以不用是 GOPATH
git clone https://github.com/go-eagle/eagle

# 进入到下载目录
cd eagle

# 编译
make build

# 运行
./scripts/admin.sh start
```

### 方式二

使用脚手架，仅生成基本目录, 不包含pkg等部分公共模块目录

```bash
# 下载
go get github.com/go-eagle/eagle/cmd/eagle

export GO111MODULE=on
# 或者在.bashrc 或 .zshrc中加入
# source .bashrc 或 source .zshrc

# 使用
eagle new eagle-demo 
# 或者 
eagle new github.com/foo/bar
```

## 💻 常用命令

- make help 查看帮助
- make dep 下载 Go 依赖包
- make build 编译项目
- make gen-docs 生成接口文档
- make test-coverage 生成测试覆盖
- make lint 检查代码规范

## 🏂 模块

## 公共模块

- 图片上传(支持本地、七牛)
- 短信验证码(支持七牛)

### 用户模块

- 注册
- 登录(邮箱登录，手机登录)
- 发送手机验证码(使用七牛云服务)
- 更新用户信息
- 关注/取消关注
- 关注列表
- 粉丝列表

## 📝 接口文档

`http://localhost:8080/swagger/index.html`

## 开发规范

遵循: [Uber Go 语言编码规范](https://github.com/uber-go/guide/blob/master/style.md)

## 📖 开发规约

- [配置说明](https://github.com/go-eagle/eagle/blob/master/config)
- [错误码设计](https://github.com/go-eagle/eagle/tree/master/pkg/errcode)
- [service 的使用规则](https://github.com/go-eagle/eagle/blob/master/internal/service)
- [repository 的使用规则](https://github.com/go-eagle/eagle/blob/master/internal/repository)
- [cache 使用说明](https://github.com/go-eagle/eagle/blob/master/pkg/cache)

## 🚀 部署

### 单独部署

上传到服务器后，直接运行命令即可

```bash
./scripts/admin.sh start
```

### Docker 部署

如果安装了 Docker 可以通过下面命令启动应用：

```bash
# 运行
docker-compose up -d

# 验证
http://127.0.0.1/health
```

### Supervisord

编译并生成二进制文件

```bash
go build -o bin_eagle
```

如果应用有多台机器，可以在编译机器进行编译，然后使用rsync同步到对应的业务应用服务器

> 以下内容可以整理为脚本

```bash
export GOROOT=/usr/local/go1.13.8
export GOPATH=/data/build/test/src
export GO111MODULE=on
cd /data/build/test/src/github.com/go-eagle/eagle
/usr/local/go1.13.8/bin/go build -o /data/build/bin/bin_eagle -mod vendor main.go
rsync -av /data/build/bin/ x.x.x.x:/home/go/eagle
supervisorctl restart eagle
```

这里日志目录设定为 `/data/log`
如果安装了 Supervisord，可以在配置文件中添加下面内容(默认：`/etc/supervisor/supervisord.conf`)：

```ini
[program:eagle]
# environment=
directory=/home/go/eagle
command=/home/go/eagle/bin_eagle
autostart=true
autorestart=true
user=root
stdout_logfile=/data/log/eagle_std.log
startsecs = 2
startretries = 2
stdout_logfile_maxbytes=10MB
stdout_logfile_backups=10
stderr_logfile=/data/log/eagle_err.log
stderr_logfile_maxbytes=10MB
stderr_logfile_backups=10
```

重启 Supervisord

```bash
supervisorctl restart eagle
```

## 📜 CHANGELOG

- [更新日志](https://github.com/go-eagle/eagle/blob/master/CHANGELOG.md)

## 🏘️ 谁在用

- [1024课堂](https://1024casts.com)
- [FastIM](https://github.com/1024casts/fastim)
- [Go微服务实战项目](https://github.com/go-microservice)

## 💬 Discussion

- Issue: https://github.com/go-eagle/eagle/issues
- QQ交流群：1074476202
- Gitter: https://gitter.im/go-eagle/eagle
- 微信交流群
<img src="https://user-images.githubusercontent.com/3043638/159420999-e00a667d-a5d9-404b-876a-ba0bc94981b9.jpeg" width="200px">

## Stargazers over time

[![Stargazers over time](https://starchart.cc/go-eagle/eagle.svg)](https://starchart.cc/go-eagle/eagle)

## 🔋 JetBrains 开源证书支持

`eagle` 项目一直以来都是在 JetBrains 公司旗下的 GoLand 集成开发环境中进行开发，基于 **free JetBrains Open Source license(s)** 正版免费授权，在此表达我的谢意。

<a href="https://www.jetbrains.com/?from=go-eagle/eagle" target="_blank"><img src="https://raw.githubusercontent.com/panjf2000/illustrations/master/jetbrains/jetbrains-variant-4.png" width="200" align="middle"/></a>

## 📄 License

MIT. See the [LICENSE](LICENSE) file for details.
