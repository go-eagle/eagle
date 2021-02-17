
# 分布式链路追踪

主要是使用 opentracing 协议，基于 jaeger client 来使用

## 本地快速部署

```
docker run -d -p 6831:6831/udp -p 16686:16686 jaegertracing/all-in-one:latest
```

## Reference

- https://opentracing.io/guides/golang/quick-start/
- https://www.jaegertracing.io/
- https://github.com/jaegertracing/
- https://logz.io/blog/go-instrumentation-distributed-tracing-jaeger/
- https://github.com/albertteoh/jaeger-go-example
- https://github.com/go-gorm/opentracing