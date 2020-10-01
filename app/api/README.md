# handler

接口实现层，可以理解成 MVC 的控制器层。主要接收参数、验证参数、调用service层的业务逻辑处理，最后返回数据。

PS: 如果需要进行转换数据，可以调用对应的idl进行统一数据转换。

## API风格和媒体类型说明

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

## 单元测试

接口的测试可以参考这几篇文章
 - https://rshipp.com/go-api-integration-testing/
 - https://github.com/quii/learn-go-with-tests