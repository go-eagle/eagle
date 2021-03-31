
# 分布式链路追踪

主要是使用 opentracing 协议，基于 jaeger client 来使用

## 本地快速部署

```
docker run -d -p 6831:6831/udp -p 16686:16686 jaegertracing/all-in-one:latest
```


## 主要使用步骤

### 1、初始化 Jaeger 并将 tracer设为全局，方便后续调用

```
// initJaeger 将jaeger tracer设置为全局tracer
func initJaeger(service string) (opentracing.Tracer, io.Closer) {
	cfg := jaegercfg.Configuration{
		// 将采样频率设置为1，每一个span都记录，方便查看测试结果
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
			// 将span发往jaeger-collector的服务地址
			CollectorEndpoint: "http://localhost:14268/api/traces",
		},
	}

    tracer, closer, err := cfg.NewTracer(
		config.Logger(jaeger.StdLogger),
		config.ZipkinSharedRPCSpan(true),
	)
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}
```

### 2、在 main函数或者中间件中获取全局 tracer，创建 root span 并执行我们第一个服务(MyFirstSpan)

```
func main() {
	closer := initJaeger("MyProcess")
	defer closer.Close()
	// 获取jaeger tracer
	t := opentracing.GlobalTracer()
	// 创建root span
	sp := t.StartSpan("MyServer")
	// main执行完结束这个span
	defer sp.Finish()
	// 将span传递给MyFirstSpan
	ctx := opentracing.ContextWithSpan(context.Background(), sp)
	MyFirstSpan(ctx)
}
```

### 3、在 MyFirstSpan 中调用另一个服务(MySecondSpan)

```
func MyFirstSpan(ctx context.Context) {
	// 开始一个span, 设置span的operation_name=MyFirstSpan
	span, ctx := opentracing.StartSpanFromContext(ctx, "MyFirstSpan")
	defer span.Finish()
	// 将context传递给MySecondSpan
	MySecondSpan(ctx)
	// 假设执行了 http请求
	span.SetTag("http_method", "GET")
	span.SetTag("http_code", "200")
	// 模拟执行耗时
	time.Sleep(1 * time.Second)
}
```

### 4、MySecondSpan:

```
func MySecondSpan(ctx context.Context) {
	// 开始一个span，设置span的operation_name=MySecondSpan
	span, ctx := opentracing.StartSpanFromContext(ctx, "MySecondSpan")
	defer span.Finish()
	// 模拟执行耗时
	time.Sleep(2 * time.Second)
	// 假设MySecondSpan发生了某些错误
	err := errors.New("something wrong")
	// 记录 log
	span.LogFields(
		log.String("event", "error"),
		log.String("message", err.Error()),
		log.Int64("error time ", time.Now().Unix()),
	)
	span.SetTag("error", true)
}
```

### 5、执行完后在 127.0.0.1:16686 查看本次 tracer

```
![tracer demo](https://static001.geekbang.org/infoq/14/14e851a82437efd4b812375558f7178b.png)
```

## 使用demo

### go-redis 

```
    span := opentracing.SpanFromContext(c.Request.Context())
	if span == nil {
		api.SendResponse(c, errno.ErrUserNotFound, "span is nil")
		return
	}
	err := tracing.Inject(span, c.Request)
	if err != nil {
		api.SendResponse(c, errno.ErrUserNotFound, err.Error())
		return
	}

	_ = apmgoredis.Wrap(redis.RedisClient).WithContext(c.Request.Context()).Set("test", 1, 100000).Err()
	_ = apmgoredis.Wrap(redis.RedisClient).WithContext(c.Request.Context()).Get("test").Err()
```

## FAQ

1、如何生成span

```
var sp opentracing.Span
carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
ctx, _ := tracer.Extract(opentracing.HTTPHeaders, carrier)
sp = tracer.StartSpan(c.Request.URL.Path, ext.RPCServerOption(ctx))
defer sp.Finish()
```

2、如何注册span到下游

```
// 生成span
span := opentracing.SpanFromContext(c.Request.Context())
if span == nil {
    api.SendResponse(c, errno.ErrUserNotFound, "span is nil")
    return
}
// 注入
err := tracing.Inject(span, c.Request)
if err != nil {
    api.SendResponse(c, errno.ErrUserNotFound, err.Error())
    return
}
```

3、如何记录tag

```
// record HTTP method
ext.HTTPMethod.Set(sp, c.Request.Method)
// record HTTP url
ext.HTTPUrl.Set(sp, c.Request.URL.String())
...
```

4、如果记录日志到span

参考：github.com/opentracing/opentracing-go/span.go

```
span.LogFields(
       log.String("event", "soft error"),
       log.String("type", "cache timeout"),
      log.Int("waited.millis", 1500))
```

## Reference

- https://opentracing.io/guides/golang/quick-start/
- https://www.jaegertracing.io/
- https://github.com/jaegertracing/
- https://medium.com/opentracing/take-opentracing-for-a-hotrod-ride-f6e3141f7941
- demo project: https://github.com/jaegertracing/jaeger/tree/master/examples/hotrod
- https://github.com/opentracing-contrib?q=&type=&language=go
- https://logz.io/blog/go-instrumentation-distributed-tracing-jaeger/
- https://github.com/albertteoh/jaeger-go-example
- https://github.com/go-gorm/opentracing
- https://github.com/opentracing-contrib/go-gin/blob/master/ginhttp/server.go
- https://github.com/opentracing-contrib/go-gin/blob/master/examples/example_test.go
- https://github.com/opentracing-contrib/goredis
- https://xie.infoq.cn/article/6450b96c33298bab92ba6f3c2
- [如何在Go中使用OpenTelemetry开始分布式追踪](https://www.youtube.com/watch?v=yQpyIrdxmQc)