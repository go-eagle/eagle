
通过启动 prometheus 来采集服务的各种运行指标，方便监控。

启动一个 prometheus

```bash
prometheus --config.file=prometheus.yml
```

可以通过 http://localhost:8090/ping  来检测服务的运行情况


## Reference

- https://github.com/yolossn/Prometheus-Basics