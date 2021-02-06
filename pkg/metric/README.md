
监控主要分为3部分

- 请求监控基本指标: QPS、请求响应时间、请求体大小和响应体大小，可以使用本框架的 `PromMiddleware` 中间件来进行监控
- Go 进程监控: Goroutine使用数量、GC回收时间、内存使用情况、活动对象数、对象分配的速率等，框架已内置 `/metrics`, 只需要在 `Prometheus` 接入即可。 可以使用插件6671来查看
- 具体业务监控: 比如cache miss等

Go 进程监控图示例
![Go processes](https://grafana.com/api/dashboards/6671/images/4286/image)

应用监控 [面板配置](golang_app_dashboard.json) ，可以直接导入到Grafana进行使用。