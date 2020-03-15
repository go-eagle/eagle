# 说明

本logger主要是对zap库的封装，便于使用。  

## 日志功能
- 将日志信息记录到日志文件里
- 日志切割-能够根据日志文件大小或时间间隔进行切割
- 支付不同的日志级别(eg：info,debug,warn,error,fatal)
- 能够打印基本信息，如调用文件/函数名和行号，日志时间，IP等

## 使用方法


## Reference
 - 日志基础库zap: https://github.com/uber-go/zap
 - 日志分割库-按时间：https://github.com/lestrrat-go/file-rotatelogs
 - 日志分割库-按大小：https://github.com/natefinch/lumberjack 
 - [深度 | 从Go高性能日志库zap看如何实现高性能Go组件](https://mp.weixin.qq.com/s/i0bMh_gLLrdnhAEWlF-xDw)
 - https://github.com/wzyonggege/logger
 - https://wisp888.github.io/golang-iris-%E5%AD%A6%E4%B9%A0%E4%BA%94-zap%E6%97%A5%E5%BF%97.html
