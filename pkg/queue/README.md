## 消息队列

- RocketMQ
- RabbitMQ
- Kafka
- Nats


## 作用

- 系统解耦
- 异步处理
- 削峰填谷

## 客户端

- RocketMQ Go客户端: https://github.com/apache/rocketmq-client-go
- kafka-go: https://github.com/segmentio/kafka-go
- Nats.go: github.com/nats-io/nats.go

> 如果是阿里云RocketMQ: 可以使用官方自己出的库


## 注意事项

- Consumer 因可能多次收到同一消息，需要做好幂等处理
- 消费时记录日志，方便后续定位问题，最好加上请求的唯一标识，比如 request_id或trace_id之类的字段
- 尽量使用批量方式消费，可以很大程度上提高消费吞吐量


## Reference

- [RocketMQ官网](https://rocketmq.apache.org/)
- [RocketMQ文档](https://rocketmq.apache.org/docs/quick-start/)
- [RocketMQ Go客户端](https://github.com/apache/rocketmq-client-go)
- [RocketMQ Go客户端使用文档](https://github.com/apache/rocketmq-client-go/blob/master/docs/Introduction.md)
- [阿里云RocketMQ](https://cn.aliyun.com/product/rocketmq)
- https://github.com/GSabadini/go-message-broker/blob/master/main.go
- [Automatically recovering RabbitMQ connections in Go applications](https://medium.com/@dhanushgopinath/automatically-recovering-rabbitmq-connections-in-go-applications-7795a605ca59)