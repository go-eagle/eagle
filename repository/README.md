# repository

数据访问层，负责访问 DB、MC、外部 HTTP 等接口，对上层屏蔽数据访问细节。

具体职责有：

SQL 拼接和 DB 访问逻辑
DB 的拆库折表逻辑
DB 的缓存读写逻辑
HTTP 接口调用逻辑

参考：
https://github.com/realsangil/apimonitor/blob/fe1e9ef75dfbf021822d57ee242089167582934a/pkg/rsdb/repository.go