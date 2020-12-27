
# 限流

限流分为单机限流和分布式限流。

单机限流一般分为 令牌桶算法和漏斗算法限流。

## 限流算法

### 令牌桶算法限流

代码实现：[golang.org/x/time/rate](golang.org/x/time/rate)

### 漏斗算法限流

代码实现：[https://github.com/uber-go/ratelimit](https://github.com/uber-go/ratelimit)

## 过载保护

但是采用漏斗桶/令牌桶的缺点是太被动, 不能快速适应流量变化。  
因此我们需要一种自适应的限流算法，即: 过载保护，根据系统当前的负载自动丢弃流量。

可以参考 BBR算法实现： [https://github.com/go-kratos/kratos/blob/master/pkg/ratelimit/bbr/bbr.go](https://github.com/go-kratos/kratos/blob/master/pkg/ratelimit/bbr/bbr.go)

### Sentinel-golang 版本 

sentinel 官网：https://sentinelguard.io/zh-cn/index.html

sentinel 是面向分布式服务架构的高可用流量控制组件

#### Github

https://github.com/alibaba/sentinel-golang

#### 特性

Sentinel Go 1.0 版本对齐了 Java 版本核心的高可用防护和容错能力，包括

- 限流
- 流量整形
- 并发控制
- 熔断降级
- 系统自适应保护
- 热点防护

等特性

更多使用介绍看这里：[流控降级组件 Sentinel Go简介](https://mp.weixin.qq.com/s/j1kTArkROXlymghR1hkDFA)

## Reference

- [Sentinel Go版本](https://github.com/alibaba/sentinel-golang)
- [常用微服务框架的适配](https://github.com/sentinel-group/sentinel-go-adapters)
- [动态数据源扩展支持](https://github.com/sentinel-group/sentinel-go-datasources)
- [阿里双11同款流控降级组件 Sentinel Go简介](https://mp.weixin.qq.com/s/j1kTArkROXlymghR1hkDFA)
- [分布式高并发服务三种常用限流方案简介](https://mp.weixin.qq.com/s/zIhQuK1jmHcn5eIqhJfNkw)
