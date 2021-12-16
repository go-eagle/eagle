# conf

默认从 config.local.yaml 加载配置, 可以通过下面命令生成：

```bash
cp config.sample.yaml config.local.yaml
```

如果是多环境可以设定不同的配置文件，比如：
 - config.local.yaml 本地开发环境
 - config.dev.yaml 开发环境
 - config.test.yaml 测试环境
 - config.pre.yaml 预发布环境
 - config.prod.yaml 线上环境
 
 生成开发环境配置文件：
 
 ```bash
cp config.sample.yaml config.dev.yaml
```
 
 以上配置是可以提交到Git等代码仓库的，`config.local.yaml` 配置文件默认不会被提交。
 
 使用本地配置文件来运行程序，命令如下:
 
 ```bash
# 本地启动
# 也可以直接 ./eagle
./eagle -c config/config.local.yaml
# 开发环境启动
./eagle -c config/config.dev.yaml
# 线上环境启动
./eagle -c config/config.prod.yaml

```
