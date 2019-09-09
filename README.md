snake 框架

适合于开发Go项目,主要是提供API 服务。

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

> 参考 新浪开放平台 [Error code](http://open.weibo.com/wiki/Error_code) 的设计

错误返回值格式：

```json
{
  "code": 10002,
  "message": "Error occurred while binding the request body to the struct."
}
```
<details>
<summary>错误代码说明：</summary>

| 1 | 00 | 02 |
| :------ | :------ | :------ |
| 服务级错误（1为系统级错误） | 服务模块代码 | 具体错误代码 |

- 服务级别错误：1 为系统级错误；2 为普通错误，通常是由用户非法操作引起的
- 服务模块为两位数：一个大型系统的服务模块通常不超过两位数，如果超过，说明这个系统该拆分了
- 错误码为两位数：防止一个模块定制过多的错误码，后期不好维护
- `code = 0` 说明是正确返回，`code > 0` 说明是错误返回
- 错误通常包括系统级错误码和服务级错误码
- 建议代码中按服务模块将错误分类
- 错误码均为 >= 0 的数
- 在本项目中 HTTP Code 固定为 http.StatusOK，错误码通过 code 来表示。
</details>

## 模块

好了，现在可以使用该项目来开发相关业务了。

Happy Coding. ^_^
