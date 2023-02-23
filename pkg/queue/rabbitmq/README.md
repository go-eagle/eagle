# RabbitMQ

## 什么是RabbitMQ？

RabbitMQ是一套开源（MPL）的消息队列服务软件，是由 LShift 提供的一个 Advanced Message Queuing Protocol (AMQP) 的开源实现，由以高性能、健壮以及可伸缩性出名的 Erlang 写成。

RabbitMQ的特点：

- **可靠性**。支持持久化，传输确认，发布确认等保证了MQ的可靠性。
- **灵活的分发消息策略**。这应该是RabbitMQ的一大特点。在消息进入MQ前由Exchange(交换机)进行路由消息。分发消息策略有：简单模式、工作队列模式、发布订阅模式、路由模式、通配符模式。
- **支持集群**。多台RabbitMQ服务器可以组成一个集群，形成一个逻辑Broker。
- **多种协议**。RabbitMQ支持多种消息队列协议，比如 STOMP、MQTT 等等。
- **支持多种语言客户端**。RabbitMQ几乎支持所有常用编程语言，包括 Java、.NET、Ruby 等等。
- **可视化管理界面**。RabbitMQ提供了一个易用的用户界面，使得用户可以监控和管理消息 Broker。
- **插件机制**。RabbitMQ提供了许多插件，可以通过插件进行扩展，也可以编写自己的插件。


## AMQP基础概念

AMQP是一套公开的消息队列协议，最早在2003年被提出，它旨在从协议层定义消息通信数据的标准格式，为的就是解决MQ市场上协议不统一的问题。RabbitMQ就是遵循AMQP标准协议开发的MQ服务。

- 即Advanced Message Queuing Protocol，一个提供统一消息服务的应用层标准高级消息队列协议,是应用层协议的一个开放标准，为面向消息的中间件设计；
- AMQP 的主要特征是面向消息、队列、路由（包括点对点和发布/订阅）、可靠性、安全。
- RabbitMQ 是一个开源的 AMQP 实现，服务器端用Erlang语言编写，支持多种客户端，如：Python、Ruby、.NET、Java、PHP等。

### Producer（生产者）

消息生产者。

从安全角度考虑，网络是不可靠的，接收消息的应用也有可能在处理消息的时候失败。基于此原因，AMQP模块包含了一个消息确认（message acknowledgements）的概念：当一个消息从队列中投递给消费者后（Consumer），消费者会通知一下消息代理（Broker），这个可以是自动的，也可以由处理消息的应用的开发者执行。当“消息确认”被启用的时候，消息代理不会完全将消息从队列中删除，直到它收到来自消费者的确认回执（acknowledgement）。

### Consumer（消费者）

消息消费者。

### Connection（连接）

一个网络连接，比如TCP/IP套接字连接。Channel是建立在Connection之上的，一个Connection可以建立多个Channel。

### Channel（信道）

信道是多路复用连接中的一条独立的双向数据流通道，为会话提供物理传输介质。Channel是在connection内部建立的逻辑连接，如果应用程序支持多线程，通常每个thread创建单独的channel进行通讯，AMQP method包含了channel id帮助客户端和message broker识别channel，所以channel之间是完全隔离的。Channel作为轻量级的Connection极大减少了操作系统建立TCP connection的开销。在客户端的每个连接里，可建立多个Channel，每个Channel代表一个会话任务。

### Broker（消息代理）

其实Broker就是接收和分发消息的应用，也就是说RabbitMQ Server就是Message Broker。

### Vhost（虚拟主机）

虚拟主机，，一批交换器（Exchange），消息队列（Queue）和相关对象。虚拟主机是共享相同身份认证和加密环境的独立服务器域。同时一个Broker里可以开设多个vhost，用作不同用户的权限分离。

### Exchange（交换机）

在RabbitMQ中，生产者将消息发送到Exchange，而不是队列（Queue）之中。消息是由Exchange路由到一个或多个队列之中，如果路由不到，或返回给生产者、或直接丢弃。

#### 交换机的类型

Exchange有4种类型对应4种不同的路由策略:

##### 1. Fanout（扇型交换机）

针对队列的广播，它会忽略BindingKey，将所有发送到该Exchange的消息路由到所有与该Exchange绑定的队列中。

##### 2. Direct（直连交换机）

它会将消息路由到那些RoutingKey和BindingKey完全一样的队列中。

##### 3. Topic（主题交换机）

与direct类似，只不过不要求RoutingKey和BindingKey完全一致，可以模糊匹配。

##### 4. Headers（头交换机）

根据消息内容中的headers属性进行匹配，很少用。

#### 交换机的状态

交换机可以有两个状态：

- 持久（durable）
- 暂存（transient）

#### 交换机的属性

- Name
- Durability （消息代理重启后，交换机是否还存在）
- Auto-delete （当所有与之绑定的消息队列都完成了对此交换机的使用后，删掉它）
- Arguments（依赖代理本身）

### Queue（消息队列）

是 RabbitMQ 的内部对象，用于存储消息。每个消息都会被投入到一个或多个队列。且多个消费者可以订阅同一个 Queue（这时 Queue 中的消息会被平均分摊给多个消费者进行处理，而不是每个消费者都收到所有的消息并处理）。

#### 属性

- Name
- Durable（消息代理重启后，队列依旧存在）
- Exclusive（只被一个连接（connection）使用，而且当连接关闭后队列即被删除）
- Auto-delete（当最后一个消费者退订后即被删除）
- Arguments（一些消息代理用他来完成类似与TTL的某些额外功能）

### Binding（绑定）

它的作用就是把Exchange（Exchange）和队列（Queue）关联起来，在绑定的时候一般会指定一个BindingKey。

### Routing Key（路由键）

生产者将消息发送给Exchange时，一般会指定一个RoutingKey，Exchange会根据这个值选择一些路由规则。

### Binding Key（绑定键）

指定当前 Exchange（交换机）下，什么样的 Routing Key（路由键）会被下派到当前绑定的 Queue 中。

## Docker部署开发环境

```shell
docker pull bitnami/rabbitmq:latest

docker run -itd \
    --hostname localhost \
    --name rabbitmq-test \
    -p 15672:15672 \
    -p 5672:5672 \
    -p 1883:1883 \
    -p 15675:15675 \
    -e RABBITMQ_PLUGINS=rabbitmq_top,rabbitmq_mqtt,rabbitmq_web_mqtt,rabbitmq_prometheus,rabbitmq_stomp,rabbitmq_auth_backend_http \
    bitnami/rabbitmq:latest

# 查看插件列表
rabbitmq-plugins list
# rabbitmq_peer_discovery_consul 
rabbitmq-plugins --offline enable rabbitmq_peer_discovery_consul
# rabbitmq_mqtt 提供与后端服务交互使用，端口1883
rabbitmq-plugins enable rabbitmq_mqtt
# rabbitmq_web_mqtt 提供与前端交互使用，端口15675
rabbitmq-plugins enable rabbitmq_web_mqtt
# 
rabbitmq-plugins enable rabbitmq_auth_backend_http
```

管理后台: <http://localhost:15672>  
默认账号: guest  
默认密码: guest

## Reference

- https://github.com/rabbitmq/amqp091-go
- [https://github.com/rabbitmq/amqp091-go/tree/v1.6.1/_examples](https://github.com/rabbitmq/amqp091-go/tree/main/_examples)
- https://medium.com/@dhanushgopinath/automatically-recovering-rabbitmq-connections-in-go-applications-7795a605ca59
- https://blog.boot.dev/golang/connecting-to-rabbitmq-in-golang-easy/
- https://github.com/wagslane/go-rabbitmq
