# conf

默认从 config/{ENV} 加载配置, 可以通过下面命令生成：

如果是多环境可以设定不同的配置目录，比如：
 - config/local 本地开发环境
 - config/test 测试环境
 - config/staging 预发布环境
 - config/prod 线上环境
 
 > 环境命名参考自：https://cloud.google.com/apis/design/directory_structure
 
 生成测试环境配置文件：
 
 ```bash
cp -r config/localhost config/test
```
 
 使用本地配置文件来运行程序，命令如下:
 
 ```bash
# 本地启动
# 也可以直接 APP_ENV={env} ./eagle
APP_ENV=local ./eagle -c config
或
./eagle -e local -c config
```
