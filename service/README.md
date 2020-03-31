# service

 - 业务逻辑层，处于 handler 层和 repository 层之间 
 - service 只能通过 repository 层获取数据
 
一个业务一个目录，比如用户是在user目录下，设计用户相关的都可以放到这里，根据不同的模块分离到不同的文件，同时又避免了单个文件func太多的问题。
比如：
 - 用户基础服务- user_service.go
 - 用户关注- user_follow_service.go
 - 用户喜欢- user_like_service.go
 - 用户评论- user_comment_service.go

 
 ## Reference
 - https://github.com/qiangxue/go-rest-api
 - https://github.com/irahardianto/service-pattern-go
 - https://github.com/golang-standards/project-layout
 - https://www.youtube.com/watch?v=oL6JBUk6tj0
 - https://peter.bourgon.org/blog/2017/06/09/theory-of-modern-go.html
 - https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1