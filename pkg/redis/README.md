
单元测试可以使用 https://github.com/alicebob/miniredis, 可以开启一个本地的模拟redis

- [在单元测试中模拟Redis](https://medium.com/@elliotchance/mocking-redis-in-unit-tests-in-go-28aff285b98)

## 案例

- [Redis分布式锁没用明白，搞出了大故障…](https://mp.weixin.qq.com/s/BO-gly5iGLVmuG5B_FIpoQ)
- [看完这篇Redis缓存三大问题](https://mp.weixin.qq.com/s/HjzwefprYSGraU1aJcJ25g)

## Redis 优化方向

### 参数优化

maxIdle设置高点，可以保证突发流量情况下，能够有足够的连接去获取redis，不用在高流量情况下建立连接

**go-redis参数优化**

```yaml
  min_idle_conn: 30               
  dial_timeout: "1s"
  read_timeout: "500ms"
  write_timeout: "500ms"
  pool_size: 500
  pool_timeout: "60s"
```

**redisgo参数优化**

```yaml
maxIdle = 30
maxActive = 500
dialTimeout = "1s"
readTimeout = "500ms"
writeTimeout = "500ms"
idleTimeout = "60s"
```

### 使用优化

- 增加redis从库
- 对批量数据，根据redis从库数量，并发goroutine拉取数据
- 对批量数据大量使用pipeline指令
- 精简key字段
- redis的value存储编解码改为msgpack

## Pipeline
- https://redis.io/topics/pipelining
- [兼容go redis cluster的pipeline批量](http://xiaorui.cc/archives/5557)
- https://www.tizi365.com/archives/309.html