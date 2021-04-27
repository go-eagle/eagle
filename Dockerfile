# Compile stage
FROM golang:1.14-alpine AS builder

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
RUN apk add --no-cache git make ca-certificates tzdata \
    --repository http://mirrors.aliyun.com/alpine/v3.11/community \
    --repository http://mirrors.aliyun.com/alpine/v3.11/main

# 镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY="https://goproxy.cn,direct" \
    TZ=Asia/Shanghai

# 移动到工作目录
WORKDIR /go/src/github.com/1024casts/snake

# 复制项目中的 go.mod 和 go.sum文件并下载依赖信息
COPY go.mod .
COPY go.sum .
RUN go mod download

# 将代码复制到容器中
COPY . .
COPY config ./conf

# Build the Go app
RUN make build

# 创建一个小镜像
# Final stage
FROM debian:stretch-slim

WORKDIR /app

# 从builder镜像中把 /build 拷贝到当前目录
COPY --from=builder /go/src/github.com/1024casts/snake/snake  /app
COPY --from=builder /go/src/github.com/1024casts/snake/conf   /app

RUN mkdir -p /data/logs/

# Expose port 8080 to the outside world
EXPOSE 8080

# 需要运行的命令
CMD ["/app/snake", "-c", "conf/config.docker.yaml"]

# 1. build image: docker build -t snake:v1 -f Dockerfile .
# 2. start: docker run --rm -it -p 8080:8080 snake:v1
# 3. test: curl -i http://localhost:8080/health

