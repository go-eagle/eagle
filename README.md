# 🐍 snake

 [![GitHub Workflow Status](https://img.shields.io/github/workflow/status/1024casts/snake/Go?style=flat-square)](https://github.com/1024casts/snake)
 [![codecov](https://codecov.io/gh/1024casts/snake/branch/master/graph/badge.svg)](https://codecov.io/gh/1024casts/snake)
 [![GolangCI](https://golangci.com/badges/github.com/golangci/golangci-lint.svg)](https://golangci.com)
 [![godoc](https://godoc.org/github.com/1024casts/snake?status.svg)](https://godoc.org/github.com/1024casts/snake)
 [![Go Report Card](https://goreportcard.com/badge/github.com/1024casts/snake)](https://goreportcard.com/report/github.com/1024casts/snake)
 [![gitmoji](https://img.shields.io/badge/gitmoji-%20%F0%9F%98%9C%20%F0%9F%98%8D-FFDD67.svg?style=flat-square)](https://github.com/carloscuesta/gitmoji)
 [![License](https://img.shields.io/github/license/1024casts/snake?style=flat-square)](/LICENSE)

一款适合于快速开发业务的 Go 框架，可快速构建 API 服务 或 Web 网站。

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
- 测试框架 [GoConvey](http://goconvey.co/)
- CI/CD [GitHub Actions](https://github.com/actions)
- 使用 [GolangCI-lint](https://golangci.com/) 进行代码检测
- 使用 make 来管理 Go 工程
- 使用 shell(admin.sh) 脚本来管理进程
- 使用 YAML 文件进行多环境配置

## 📗 目录结构

```shell
├── Makefile                     # 项目管理文件
├── build                        # 编译目录
├── cmd                          # 脚手架目录
├── conf                         # 配置文件统一存放目录
├── db.sql                       # 在部署新环境时，可以登录 MySQL 客户端，执行 source db.sql 创建数据库和表
├── docs                         # Swagger 文档，执行 swag init 生成的
├── handler                      # 控制器目录，用来读取输入、调用业务处理、返回结果
├── internal                     # 业务目录
│   ├── cache                    # 基于业务封装的cache
│   ├── idl                      # 数据结构转换
│   ├── model                    # 数据库 model
│   ├── repository               # 数据访问层
│   └── service                  # 业务逻辑层
├── logs                         # 存放日志的目录
├── main.go                      # 项目入口文件
├── pkg                          # 一些封装好的 package
├── router                       # 路由及中间件目录
└── scripts                      # 存放用于执行各种构建，安装，分析等操作的脚本
```

## 🛠️ 快速开始

### 方式一

直接Clone项目的方式，文件比较全

TIPS: 需要本地安装MySQL数据库和 Redis

```bash
# 下载安装，可以不用是 GOPATH
git clone https://github.com/1024casts/snake

# 进入到下载目录
cd snake


# 生成本地环境配置文件
cd conf
cp config.sample.yaml config.local.yaml

# 编译
make build

# 运行
./scripts/admin.sh start
```

### 方式二

使用脚手架，仅生成基本目录, 不包含pkg等部分公共模块目录

```bash
# 下载
go get github.com/1024casts/snake/cmd/snake

cd $GOPATH/src
# 使用
snake new snake-demo -d ./
# 或者 
snake new github.com/foo/bar -d ./
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

遵循: [Uber Go 语言编码规范](https://github.com/xxjwxc/uber_go_guide_cn)

## 📖 开发规约

- [配置说明](https://github.com/1024casts/snake/blob/master/conf)
- [错误码设计](https://github.com/1024casts/snake/tree/master/pkg/errno)
- [service 的使用规则](https://github.com/1024casts/snake/blob/master/internal/service)
- [repository 的使用规则](https://github.com/1024casts/snake/blob/master/internal/repository)
- [cache 使用说明](https://github.com/1024casts/snake/blob/master/pkg/cache)

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
go build -o bin_snake
```

如果应用有多台机器，可以在编译机器进行编译，然后使用rsync同步到对应的业务应用服务器

> 以下内容可以整理为脚本

```bash
export GOROOT=/usr/local/go1.13.8
export GOPATH=/data/build/test/src
export GO111MODULE=on
cd /data/build/test/src/github.com/1024casts/snake
/usr/local/go1.13.8/bin/go build -o /data/build/bin/bin_snake -mod vendor main.go
rsync -av /data/build/bin/ x.x.x.x:/home/go/snake
supervisorctl restart snake
```

这里日志目录设定为 `/data/log`
如果安装了 Supervisord，可以在配置文件中添加下面内容(默认：`/etc/supervisor/supervisord.conf`)：

```ini
[program:snake]
# environment=
directory=/home/go/snake
command=/home/go/snake/bin_snake
autostart=true
autorestart=true
user=root
stdout_logfile=/data/log/snake_std.log
startsecs = 2
startretries = 2
stdout_logfile_maxbytes=10MB
stdout_logfile_backups=10
stderr_logfile=/data/log/snake_err.log
stderr_logfile_maxbytes=10MB
stderr_logfile_backups=10
```

重启 Supervisord

```bash
supervisorctl restart snake
```

## 📜 CHANGELOG

- [更新日志](https://github.com/1024casts/snake/blob/master/CHANGELOG.md)

## 🏘️ 谁在用

- [1024课堂](https://1024casts.com)
- [FastIM](https://github.com/1024casts/fastim)

## 💬 Discussion

- Issue: https://github.com/1024casts/snake/issues
- QQ交流群：1074476202

## 🔋 JetBrains 开源证书支持

`snake` 项目一直以来都是在 JetBrains 公司旗下的 GoLand 集成开发环境中进行开发，基于 **free JetBrains Open Source license(s)** 正版免费授权，在此表达我的谢意。

<a href="https://www.jetbrains.com/?from=1024casts/snake" target="_blank"><img src="https://raw.githubusercontent.com/panjf2000/illustrations/master/jetbrains/jetbrains-variant-4.png" width="200" align="middle"/></a>

## 📄 License

MIT. See the [LICENSE](LICENSE) file for details.
