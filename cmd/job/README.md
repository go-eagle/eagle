# 定时任务目录

`main.go` 是任务启动的入口文件，`demo` 目录是提供的演示程序

业务可以根据业务进行分组，比如  

```bash
- user   // 用户
- feed   // feed
- im     // im
- post   // 帖子
...
```

然后将要执行的任务添加到 `main.go` 。