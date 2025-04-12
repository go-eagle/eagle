# This document explains how to use the Kafka package for message queue functionality

## 初始化配置

### 1. 配置文件示例

首先需要在项目的 `config/{ENV}` 目录下创建 `kafka.yaml` 配置文件：

```yaml
default:   # 实例名称
  Version: "2.8.1"
  RequiredAcks: 1
  Topic: "test-topic"
  ConsumeTopic: 
    - "test-topic1"
    - "test-topic2"
  Brokers:
    - "localhost:9092"
  GroupID: "test-group"
  Partitioner: "hash"

order:    # 另一个实例配置
  Version: "2.8.1"
  RequiredAcks: 1
  Topic: "order-topic"
  ConsumeTopic: 
    - "order-topic"
  Brokers:
    - "localhost:9092"
  GroupID: "order-group"
  Partitioner: "random"
```

### 2. 代码使用示例

#### 基本使用

```go
package main

import (
    "fmt"

    "github.com/go-eagle/eagle/pkg/queue/kafka"
)

func main() {
    // 加载配置
    kafka.Load()
    defer kafka.Close()

    // 获取配置
    configs := kafka.GetConfig()
    
    // 获取指定实例的配置
    defaultCfg := configs["default"]
    fmt.Printf("Default kafka config: %+v\n", defaultCfg)
    
    orderCfg := configs["order"]
    fmt.Printf("Order kafka config: %+v\n", orderCfg)
}
```

#### 完整应用示例

```go
package main

import (
    "context"
    "log"
    
    "github.com/go-eagle/eagle/pkg/queue/kafka"
)

func main() {
    // 1. 初始化配置
    kafka.Load()
    defer kafka.Close()
    
    // 2. 获取配置信息（可选）
    configs := kafka.GetConfig()
    if len(configs) == 0 {
        log.Fatal("No kafka config found")
    }
    
    // 3. 使用配置进行消息发布
    ctx := context.Background()
    err := kafka.Publish(ctx, "default", "test-topic", "hello world")
    if err != nil {
        log.Printf("Failed to publish message: %v", err)
    }
    
    // 4. 使用配置进行消息消费
    handler := func(data []byte) error {
        log.Printf("Received message: %s", string(data))
        return nil
    }
    
    // 从默认实例消费
    go func() {
        err := kafka.ConsumePartition(ctx, "default", "test-topic", handler)
        if err != nil {
            log.Printf("Failed to consume message: %v", err)
        }
    }()
    
    // 从order实例消费
    go func() {
        err := kafka.ConsumePartition(ctx, "order", "order-topic", handler)
        if err != nil {
            log.Printf("Failed to consume message: %v", err)
        }
    }()
    
    select {}
}
```

### 3. 配置说明

#### Conf 结构体字段说明

```go
type Conf struct {
    Version      string   // Kafka 版本号
    RequiredAcks int      // 消息确认机制：0=不确认，1=leader确认，-1=所有副本确认
    Topic        string   // 默认主题
    ConsumeTopic []string // 消费主题列表
    Brokers      []string // Kafka broker地址列表
    GroupID      string   // 消费者组ID
    Partitioner  string   // 分区策略：random/roundrobin/hash
}
```

#### 重要方法说明

1. `Load()`: 加载配置并初始化 Kafka 管理器
2. `Close()`: 关闭 Kafka 连接
3. `GetConfig()`: 获取所有配置信息

### 4. 最佳实践

1. 在应用启动时调用 `Load()`
2. 在应用退出时调用 `Close()`
3. 使用多个实例时，通过不同的配置名区分
4. 根据业务需求选择合适的分区策略
5. 合理设置 `RequiredAcks` 确保消息可靠性

### 5. 注意事项

1. 配置文件必须位于项目的 `config` 目录下
2. 配置文件名必须为 `kafka.yaml`
3. 确保配置的 broker 地址可访问
4. Version 字段要与实际 Kafka 集群版本匹配
5. 记得在程序退出时调用 `Close()` 释放资源

## client 组件的使用

使用 `client` 可以快速方便的使用 `kafka` 组件，以下是 Kafka 配置初始化的使用说明和示例。

### Method Documentation

#### 1. Publish

用于发送消息到指定的 topic。

```go
func Publish(ctx context.Context, name string, topic, msg string) error
```

**Parameters:**

- `ctx`: 上下文
- `name`: kafka 实例名称
- `topic`: 主题名称
- `msg`: 消息内容

**Example:**

````go
package main

import (
    "context"
    "github.com/go-eagle/eagle/pkg/queue/kafka"
)

func main() {
    ctx := context.Background()
    err := kafka.Publish(ctx, "default", "test-topic", "hello world")
    if err != nil {
        panic(err)
    }
}
````

### 2. Consume

使用消费组模式消费消息。

```go
func Consume(ctx context.Context, name string, topics []string, handler sarama.ConsumerGroupHandler) error
```

**Parameters:**

- `ctx`: 上下文
- `name`: kafka 实例名称
- `topics`: 主题列表
- `handler`: 消费组处理器

**Example:**

```go
package main

import (
    "context"
    "github.com/Shopify/sarama"
    "github.com/go-eagle/eagle/pkg/queue/kafka"
)

type ConsumerGroupHandler struct {
    sarama.ConsumerGroupHandler
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for msg := range claim.Messages() {
        // 处理消息
        println(string(msg.Value))
        session.MarkMessage(msg, "")
    }
    return nil
}

func main() {
    ctx := context.Background()
    handler := &ConsumerGroupHandler{}
    err := kafka.Consume(ctx, "default", []string{"test-topic"}, handler)
    if err != nil {
        panic(err)
    }
}
```

### 3. ConsumePartition

消费指定 topic 的所有分区。

```go
func ConsumePartition(ctx context.Context, name, topic string, handler func([]byte) error) error
```

**Parameters:**

- `ctx`: 上下文
- `name`: kafka 实例名称
- `topic`: 主题名称
- `handler`: 消息处理函数

**Example:**

```go
package main

import (
    "context"
    "github.com/go-eagle/eagle/pkg/queue/kafka"
)

func main() {
    ctx := context.Background()
    handler := func(data []byte) error {
        println(string(data))
        return nil
    }
    
    err := kafka.ConsumePartition(ctx, "default", "test-topic", handler)
    if err != nil {
        panic(err)
    }
}
```

### 4. ConsumerByPartitionId

消费指定 topic 的指定分区。

```go
func ConsumerByPartitionId(ctx context.Context, name, topic string, partition int32, handler func([]byte) error) error
```

**Parameters:**

- `ctx`: 上下文
- `name`: kafka 实例名称
- `topic`: 主题名称
- `partition`: 分区 ID
- `handler`: 消息处理函数

**Example:**

```go
package main

import (
    "context"
    "github.com/go-eagle/eagle/pkg/queue/kafka"
)

func main() {
    ctx := context.Background()
    handler := func(data []byte) error {
        println(string(data))
        return nil
    }
    
    err := kafka.ConsumerByPartitionId(ctx, "default", "test-topic", 0, handler)
    if err != nil {
        panic(err)
    }
}
```

### 5. GetPartitionList

获取指定 topic 的分区列表。

```go
func GetPartitionList(ctx context.Context, name, topic string) ([]int32, error)
```

**Parameters:**

- `ctx`: 上下文
- `name`: kafka 实例名称
- `topic`: 主题名称

**Example:**

```go
package main

import (
    "context"
    "fmt"
    "github.com/go-eagle/eagle/pkg/queue/kafka"
)

func main() {
    ctx := context.Background()
    partitions, err := kafka.GetPartitionList(ctx, "default", "test-topic")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Topic partitions: %v\n", partitions)
}
```

### 使用场景示例

#### 1. 基本生产和消费

````go
// 生产者
err := kafka.Publish(ctx, "default", "test-topic", "message")

// 消费者（消费组模式）
handler := &ConsumerGroupHandler{}
err = kafka.Consume(ctx, "default", []string{"test-topic"}, handler)
````

#### 2. 分区消费

````go
// 消费所有分区
err := kafka.ConsumePartition(ctx, "default", "test-topic", handler)

// 消费指定分区
err = kafka.ConsumerByPartitionId(ctx, "default", "test-topic", 0, handler)

// 获取分区信息
partitions, err := kafka.GetPartitionList(ctx, "default", "test-topic")
````

#### 3. 完整流程示例

````go
package main

import (
    "context"
    "log"

    "github.com/go-eagle/eagle/pkg/queue/kafka"
)

func main() {
    ctx := context.Background()
    
    // 1. 获取分区信息
    partitions, err := kafka.GetPartitionList(ctx, "default", "test-topic")
    if err != nil {
        panic(err)
    }
    
    // 2. 发送消息
    err = kafka.Publish(ctx, "default", "test-topic", "test message")
    if err != nil {
        panic(err)
    }
    
    // 3. 消费消息
    handler := func(data []byte) error {
        log.Printf("Received message: %s", string(data))
        return nil
    }
    
    // 可以选择以下任一方式消费：
    
    // 3.1 消费组模式
    groupHandler := &ConsumerGroupHandler{}
    go kafka.Consume(ctx, "default", []string{"test-topic"}, groupHandler)
    
    // 3.2 消费所有分区
    go kafka.ConsumePartition(ctx, "default", "test-topic", handler)
    
    // 3.3 消费指定分区
    for _, partition := range partitions {
        go kafka.ConsumerByPartitionId(ctx, "default", "test-topic", partition, handler)
    }
    
    // 防止程序退出
    select {}
}
````

这些方法提供了灵活的消息生产和消费方式，可以根据实际需求选择合适的使用方式。

## 独立组件的使用

### 生产者(Producer) Example

如何创建并使用 Kafka producer:

```go
package main

import (
    "context"
    
    "github.com/go-eagle/eagle/pkg/queue/kafka"
    "github.com/go-eagle/eagle/pkg/log"
)

func main() {
    // Create producer config
    conf := &kafka.Conf{
        Brokers:      []string{"localhost:9092"},
        Version:      "2.8.1",
        RequiredAcks: 1, // 0: no response, 1: wait for leader, -1: wait for all
        Partitioner:  "hash", // options: random, roundrobin, hash
    }
    
    // Create producer
    producer, err := kafka.NewProducer(conf, log.GetLogger())
    if err != nil {
        panic(err)
    }
    defer producer.Close()
    
    // Publish message
    ctx := context.Background()
    err = producer.Publish(ctx, "test-topic", "Hello World")
    if err != nil {
        panic(err)
    }
}
```

### 消费组(Consumer Group) Example

Here's how to use consumer groups for parallel processing:

```go
package main

import (
    "context"
    
    "github.com/go-eagle/eagle/pkg/queue/kafka"
    "github.com/go-eagle/eagle/pkg/log"
)

func main() {
    // Create consumer config
    conf := &kafka.Conf{
        Brokers: []string{"localhost:9092"},
        Version: "2.8.1",
        GroupID: "my-group-id",
    }
    
    // Create consumer
    consumer, err := kafka.NewConsumer(conf, log.GetLogger())
    if err != nil {
        panic(err)
    }
    defer consumer.Stop()
    
    // Create handler
    handler := &kafka.ConsumerGroupHandler{}
    
    // Start consuming
    ctx := context.Background()
    topics := []string{"test-topic"}
    err = consumer.Consume(ctx, topics, handler)
    if err != nil {
        panic(err)
    }
}
```

### 基于分区的消费者(Partition-Based Consumer) Example

Here's how to consume messages from specific partitions:

```go
package main

import (
    "context"
    
    "github.com/go-eagle/eagle/pkg/queue/kafka"
    "github.com/go-eagle/eagle/pkg/log"
)

func main() {
    conf := &kafka.Conf{
        Brokers: []string{"localhost:9092"},
        Version: "2.8.1",
    }
    
    consumer, err := kafka.NewConsumer(conf, log.GetLogger())
    if err != nil {
        panic(err)
    }
    defer consumer.Stop()
    
    ctx := context.Background()
    
    // Consume from all partitions
    err = consumer.ConsumePartition(ctx, "test-topic", func(data []byte) error {
        // Handle message
        log.Printf("Received message: %s", string(data))
        return nil
    })
    
    // Consume from specific partition
    err = consumer.ConsumerByPartitionId(ctx, "test-topic", 0, func(data []byte) error {
        // Handle message
        log.Printf("Received message from partition 0: %s", string(data))
        return nil
    })
}
```

### Features

#### 消费模式

1. **消费组模式**
   - 多个消费者可以在同一个组内并行工作
   - 消息会自动分配给各个消费者
   - 通过 `consumer.Consume()` 使用

2. **分区模式**
   - 消费所有分区：使用 `consumer.ConsumePartition()`
   - 消费指定分区：使用 `consumer.ConsumerByPartitionId()`

#### 生产者特性

1. **异步生产**
   - 消息以异步方式生产以提高性能
   - 通过内部 goroutine 处理成功/错误

2. **分区策略**
   - Random（随机）：消息随机分配到各个分区
   - RoundRobin（轮询）：消息均匀分配到各个分区
   - Hash（哈希）：基于消息 key 的哈希值分配到分区

### Configuration Options

#### Producer Config

```go
type Conf struct {
    Brokers       []string // Kafka broker addresses
    Version       string   // Kafka version
    RequiredAcks  int16    // Required acknowledgments
    Partitioner   string   // Partitioning strategy
}
```

#### Consumer Config

```go
type Conf struct {
    Brokers  []string // Kafka broker addresses
    Version  string   // Kafka version
    GroupID  string   // Consumer group ID
}
```

### 最佳实践

1. 始终正确关闭生产者和消费者
2. 使用消费组实现并行处理
3. 处理成功/错误通道中的错误
4. 配置合适的消息确认级别
5. 根据使用场景选择适当的分区策略
6. 设置适当的错误处理和日志记录

### 注意事项

- 生产者实现为异步模式以提高性能
- 消费组的重新平衡是自动处理的
- 默认从最新偏移量开始消费消息
- 重连和错误处理已内置在实现中

## Reference

- https://github.com/asong2020/Golang_Dream/tree/master/code_demo/kafka_demo
