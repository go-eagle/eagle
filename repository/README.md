# repository

数据访问层，负责访问 DB、MC、外部 HTTP 等接口，对上层屏蔽数据访问细节。

具体职责有：

SQL 拼接和 DB 访问逻辑
DB 的拆库折表逻辑
DB 的缓存读写逻辑
HTTP 接口调用逻辑

如果是返回的列表，尽量返回map，方便service使用。

## 单元测试

关于数据库的单元测试可以使用到的几个库：
 - go-sqlmock https://github.com/DATA-DOG/go-sqlmock 主要用来和数据库的交互操作:增删改
 - GoMock https://github.com/golang/mock

## Reference
 - https://github.com/realsangil/apimonitor/blob/fe1e9ef75dfbf021822d57ee242089167582934a/pkg/rsdb/repository.go
 - https://youtu.be/twcDf_Y2gXY?t=636
 - [如何使用Sqlmock对GORM应用进行单元测试](https://1024casts.com/topics/R9re7QDaq8MnJoaXRZxdljbNA5BwoK)
 - [GoMock快速入门指南](https://medium.com/better-programming/a-gomock-quick-start-guide-71bee4b3a6f1)