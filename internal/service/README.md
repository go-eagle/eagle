# service

 - 业务逻辑层，处于 `handler` 层和 `repository` 层之间 
 - `service` 只能通过 `repository` 层获取数据
 - 面向接口编程
 - 依赖接口，不要依赖实现
 - 如果有事务处理，在这一层进行处理
 - 如果是调用的第三方服务，请不要加 `cache`, 避免缓存不一致(对方更新数据，这边无法知晓)
 - 由于 `service` 会被 `http` 或 `rpc` 调用，默认会提供 `http` 调用的，比如：`GetUserInfo()`，
   如果 `rpc` 需要调用，可以对 `GetUserInfo()` 进行一层封装, 比如：`GetUser()`。
 
 ## Reference
 
 - https://github.com/qiangxue/go-rest-api
 - https://github.com/irahardianto/service-pattern-go
 - https://github.com/golang-standards/project-layout
 - https://www.youtube.com/watch?v=oL6JBUk6tj0
 - https://peter.bourgon.org/blog/2017/06/09/theory-of-modern-go.html
 - https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1
 - [按需写 service 服务层逻辑](https://www.5-wow.com/article/detail/89)
 - [Go 编程实战：如何组织代码、编写测试？](https://www.infoq.cn/article/4TAWp8YNYcVD4t046EGd)
 - https://github.com/sdgmf/go-project-sample
 - [Golang微服务最佳实践](https://sdgmf.github.io/goproject/)
 - [layout 常见大型 Web 项目分层](https://chai2010.cn/advanced-go-programming-book/ch5-web/ch5-07-layout-of-web-project.html)
 - [在 Golang 中尝试简洁架构](https://studygolang.com/articles/12909)
