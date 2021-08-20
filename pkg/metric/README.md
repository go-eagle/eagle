
监控主要分为3部分

- 请求监控基本指标: QPS、请求响应时间、请求体大小和响应体大小，可以使用本框架的 `PromMiddleware` 中间件来进行监控
- Go 进程监控: Goroutine使用数量、GC回收时间、内存使用情况、活动对象数、对象分配的速率等，框架已内置 `/metrics`, 只需要在 `Prometheus` 接入即可。 可以使用Grafana插件6671来查看
- 具体业务监控: 比如cache miss等

Go 进程监控图示例
![Go processes](https://grafana.com/api/dashboards/6671/images/4286/image)

应用监控 [面板配置](golang_app_dashboard.json) ，可以直接导入到Grafana进行使用。

`/metrics` 访问后返回示例信息如下：

```
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 5.74e-06
go_gc_duration_seconds{quantile="0.25"} 3.1516e-05
go_gc_duration_seconds{quantile="0.5"} 0.000159558
go_gc_duration_seconds{quantile="0.75"} 0.002622183
go_gc_duration_seconds{quantile="1"} 0.002622183
go_gc_duration_seconds_sum 0.002818997
go_gc_duration_seconds_count 4
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 20
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.13.8"} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 2.1649928e+07
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 2.7186312e+07
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 1.44869e+06
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 51734
# HELP go_memstats_gc_cpu_fraction The fraction of this program's available CPU time used by the GC since the program started.
# TYPE go_memstats_gc_cpu_fraction gauge
go_memstats_gc_cpu_fraction 0.06674713461941606
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 2.38592e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 2.1649928e+07
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 4.1934848e+07
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 2.4190976e+07
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 18945
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 4.1885696e+07
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 6.6125824e+07
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 1.6127124443332071e+09
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 70679
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 13888
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 16384
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 85680
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 98304
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 3.4761504e+07
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 2.014726e+06
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 983040
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 983040
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 7.3072888e+07
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 19
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 1
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
# HELP eagle_http_request_count_total Total number of HTTP requests made.
# TYPE eagle_http_request_count_total counter
eagle_http_request_count_total{endpoint="/favicon.ico",method="GET",status="404"} 1
# HELP eagle_http_request_duration_seconds HTTP request latencies in seconds.
# TYPE eagle_http_request_duration_seconds histogram
eagle_http_request_duration_seconds_bucket{endpoint="/favicon.ico",method="GET",status="404",le="0.005"} 1
eagle_http_request_duration_seconds_bucket{endpoint="/favicon.ico",method="GET",status="404",le="0.01"} 1
eagle_http_request_duration_seconds_bucket{endpoint="/favicon.ico",method="GET",status="404",le="0.025"} 1
eagle_http_request_duration_seconds_bucket{endpoint="/favicon.ico",method="GET",status="404",le="0.05"} 1
eagle_http_request_duration_seconds_bucket{endpoint="/favicon.ico",method="GET",status="404",le="0.1"} 1
eagle_http_request_duration_seconds_bucket{endpoint="/favicon.ico",method="GET",status="404",le="0.25"} 1
eagle_http_request_duration_seconds_bucket{endpoint="/favicon.ico",method="GET",status="404",le="0.5"} 1
eagle_http_request_duration_seconds_bucket{endpoint="/favicon.ico",method="GET",status="404",le="1"} 1
eagle_http_request_duration_seconds_bucket{endpoint="/favicon.ico",method="GET",status="404",le="2.5"} 1
eagle_http_request_duration_seconds_bucket{endpoint="/favicon.ico",method="GET",status="404",le="5"} 1
eagle_http_request_duration_seconds_bucket{endpoint="/favicon.ico",method="GET",status="404",le="10"} 1
eagle_http_request_duration_seconds_bucket{endpoint="/favicon.ico",method="GET",status="404",le="+Inf"} 1
eagle_http_request_duration_seconds_sum{endpoint="/favicon.ico",method="GET",status="404"} 4.2263e-05
eagle_http_request_duration_seconds_count{endpoint="/favicon.ico",method="GET",status="404"} 1
# HELP eagle_http_request_size_bytes HTTP request sizes in bytes.
# TYPE eagle_http_request_size_bytes summary
eagle_http_request_size_bytes_sum{endpoint="/favicon.ico",method="GET",status="404"} 597
eagle_http_request_size_bytes_count{endpoint="/favicon.ico",method="GET",status="404"} 1
# HELP eagle_http_response_size_bytes HTTP request sizes in bytes.
# TYPE eagle_http_response_size_bytes summary
eagle_http_response_size_bytes_sum{endpoint="/favicon.ico",method="GET",status="404"} 19
eagle_http_response_size_bytes_count{endpoint="/favicon.ico",method="GET",status="404"} 1
# HELP eagle_uptime HTTP service uptime.
# TYPE eagle_uptime counter
eagle_uptime 9
```

## Reference

- https://github.com/rcrowley/go-metrics