# 说明

本logger主要是对zap库的封装，便于使用。当然也可以根据接口使用其他日志库，比如 `logrus`。

## 日志功能

- 将日志信息记录到日志文件里
- 日志切割-能够根据日志文件大小或时间间隔进行切割
- 支持不同的日志级别(eg：info,debug,warn,error,fatal)
- 支持按日志级别分类输出到不同日志文件
- 能够打印基本信息，如调用文件/函数名和行号，日志时间，IP等

## 使用方法

```go
log.Info("user_id is 1")
log.Warn("user is not exist")
log.Error("params error")

log.Warnf("params is empty")
...
```

## 原则

日志尽量不要在 model, repository, service中打印输出，最好使用 `errors.Wrapf` 将错误和消息返回到上层，然后在 handler 层中处理错误，
也就是通过日志打印出来。  

这样做的好处是：避免相同日志在多个地方打印，让排查问题更简单。

## Reference

 - 日志基础库 zap: https://github.com/uber-go/zap
 - 日志分割库-按时间：https://github.com/lestrrat-go/file-rotatelogs
 - 日志分割库-按大小：https://github.com/natefinch/lumberjack 
 - [深度 | 从Go高性能日志库zap看如何实现高性能Go组件](https://mp.weixin.qq.com/s/i0bMh_gLLrdnhAEWlF-xDw)
 - [Logger interface for GO with zap and logrus implementation](https://www.mountedthoughts.com/golang-logger-interface/)
 - https://github.com/wzyonggege/logger
 - https://wisp888.github.io/golang-iris-%E5%AD%A6%E4%B9%A0%E4%BA%94-zap%E6%97%A5%E5%BF%97.html
 - https://www.mountedthoughts.com/golang-logger-interface/
