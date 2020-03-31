# handler

接口实现层，可以理解成 MVC 的控制器层。主要接收参数、验证参数、调用service层的业务逻辑处理，最后返回数据。

PS: 如果需要进行转换数据，可以调用对应的idl进行统一数据转换。

## 单元测试

接口的测试可以参考这几篇文章
 - https://rshipp.com/go-api-integration-testing/
 - https://github.com/quii/learn-go-with-tests