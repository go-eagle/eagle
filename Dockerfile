FROM golang:1.13-alpine AS builder

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
RUN apk add --no-cache git=2.24.3-r0 \
    --repository http://mirrors.aliyun.com/alpine/v3.11/community \
    --repository http://mirrors.aliyun.com/alpine/v3.11/main

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 移动到工作目录：/go/src/app
WORKDIR /go/src/app

# 复制项目中的 go.mod 和 go.sum文件并下载依赖信息
COPY go.mod .
COPY go.sum .
RUN go mod download

# 将代码复制到容器中
COPY . .

# 将我们的代码编译成二进制可执行文件 snake
RUN go build -o snake .

# 接下来创建一个小镜像
FROM debian:stretch-slim

COPY wait-for.sh /
# COPY ./templates /templates
# COPY ./static /static
COPY ./conf /conf

# 从builder镜像中把 /go/src/app 拷贝到当前目录
COPY --from=builder /go/src/app/snake /

RUN set -eux; \
	apt-get update; \
	apt-get install -y \
		--no-install-recommends \
		netcat; \
        chmod 755 wait-for.sh

# Expose port 8080 to the outside world
EXPOSE 8080

# 需要运行的命令
ENTRYPOINT ["/snake", "conf/config.sample.yaml"]