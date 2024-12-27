# Changelog

## v1.10.0

- feat: support async flush log to disk
- feat: add re-entry and automatic renewal for redis lock
- feat: add gorm logger
- chore: using gorm offical plguin for tracing and metrics
- chore: optimize zap log level
- chore: improve cmd[new]

## v1.9.0

- feat: support clickhouse
- chore: upgrade go to v1.22.3
- chore: optimize log and orm

## v1.8.2
- feat: support PostgreSQL
- feat: support config multiple databases
- chore(cli): generate project by adding branch name, default is http server
- chore(db): add timeout for connect, read and write
- chore(ci): upgrade goreleaser to v2

## v1.8.1
- fix: GitHub workflow badge URL
- chore: improve cmd[run/new]
- chore: improve print log for crontab
- chore: update config in logger.yaml
- chore: upgrade go to v1.19

## v1.8.0
- feat(es): add elasticsearch client(v7)
- feat: support custom filename when init log
- feat: add crontab server and example
- refactor: improve RabbitMQ client
- chore(cli): improve gen task command
- chore(response): remove init resp and improve Error
- chore: improve redis queue and add example
- chore: improve cmd: new and run

## v1.7.0
- feat(http): can custom http status
- feat(util): add util for printing stack
- feat(grpc): add tracing and keepalive for grpc client
- feat(trace): add function trace 
- chore(auth): support ignore path for auth middleware
- chore(cli): approve gen gin and proto command
- chore(lock): approve etcd lock

## v1.6.0
- add(cli): gen model
- add(cli): add gen handler cmd
- feat: add more option for server and client

## v1.5.0
- add: cache, repo, service cmd
- add: grpc error
- optimize: proto cmd

## v1.4.0
- 新增 run,proto,protoc-gen-go-gin 命令工具
- 新增 服务注册
- 新增 examples
- 优化 replace OpenTracing with OpenTelemetry
- 优化 add trace for http client and sql

## v1.3.1
- 优化 新建项目命令 `eagle new`
- 优化 目录结构
- 新增 链路追踪 tracing 支持 jaeger、zipkin、elastic
- 新增 SQL库，支持 tracing
- 新增 工具链，可以自动生成对应的 cache 文件
- 升级 go-redis到v8.10.0版本
- 升级 gin到v1.7.2版本

## v1.3.0
- 新增 Web 路由、控制器及模板
- 新增 用户中心，包含注册、登录、粉丝列表、关注列表
- 新增 sign签名增加aes对称加密算法
- 新增 redis lock 新增了过期时间的设定
- 新增 基于prometheus的统计pkg
- 新增 errGroup pkg
- 新增 tracing pkg，基于jaeger, 方便做分布式链路追踪
- 升级 gorm升级到v2.0
- 优化 用户关系模块使用独立 service
- 优化 MySQL 和 Redis 配置从结构体中获取，替换从viper中获取
- 优化 graceful stop 方法
- 优化 token sign 方法支持自定义参数payload
- 优化 redis lock, 将token收敛到包内进行处理，减少使用的心智负担
- 优化 创建项目命令 `eagle new`
- 优化 移除cache prefix, 由用户自己定义
- 优化 数据获取，使用 singleflight，防止缓存击穿
- 优化 数据为空时缓存一分钟，防止缓存穿透

## v1.2.2
- 优化 `service` 和 `repository` 的db参数
- 优化 `swagger` 文档和 `pprof` 性能分析路由仅在 `debug` 模式下开启
- 优化 `service` 的定义及其在 `hanlder` 中的使用方式
- 新增 GRPC的Server
- 新增 `app` 目录, 将原有的 `handler` 目录迁移至 `app` 目录

## v1.2.1
- 新增 增加job目录，可以在 `cmd/job` 中定义具体任务
- 优化 http client, 支持原生(raw)包方式请求和第三方库resty
- 优化 eagle new 命令，使用新的 main.go 文件
- 优化 `service` 和 `repository` 公共方法增加首个参数 `ctx Context.context`
- 优化 修改目录 `pkg/util` 为 `pkg/utils`

## v1.2.0
- 优化 main.go 文件，使用 App 全局化配置
- 优化 将 config 目录移至到 pkg 目录，并将配置解析到结构体，方便全局直接使用
- 新增 sign 包，用于url签名校验

## v1.1.2
- 新增 脚手架，可以生成项目目录

### v1.1.1
- 新增 在internal中新增cache目录，添加用户cache,供repository层进行调用
- 新增 写入缓存时加锁，防止缓存击穿(大量请求落到db)
- 优化 优化日志配置，针对warn和error级别的错误支持打印错误堆栈
- 优化 将UserModel修改为UserBaseModel

### v1.1.0
- 新增 添加internal目录，并将idl,model,repo,service目录移至该目录
- 新增 增加常用的时间和slice函数，并分文件进行存储
- 新增 对util工具包增加单元测试
- 优化 修改部分函数名使其更加语义化

### v1.0.5
- 新增 用户模块，包含: 登录、注册、关注/取消关注、关注列表和粉丝列表、查看/更新用户信息

### v1.0.4
- 新增 cache 可以通过配置文件进行设定，支持内存、redis缓存，默认使用memory
- 更新 gin框架升级到1.6.3
- 修复 日志行数显示问题

### v1.0.3
- 新增 user_repo 增加单元测试
- 新增 支持单独、docker(docker-compose方式)、Supervisord 三种部署方式
- 优化 使用zap日志库替换原有日志库
- 移除 BaseModel，ID由各model单独管理

### v1.0.2
- 新增 健康检查接口，便于对应用进行探活检测
- 新增 metrics接口，便于prometheus进行监控
- 新增 scripts目录，将各种脚本文件统一放到此目录
- 新增 schedule目录，可以配置需要跑的计划任务
- 新增 支持发送邮件模块
- 优化 main.go
- 更新 在repo README 中增加使用建议

### v1.0.1
- 优化 service的接口命名及导出命名
- 更新 README 新增开发约定

### v1.0.0
- 新增 第一个可用版本发布
