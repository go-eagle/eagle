# Compile stage
FROM golang:1.14-alpine AS builder

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
RUN apk add update --no-cache git=2.24.3-r0 \
    apk add --no-cache ca-certificates tzdata \
    --repository http://mirrors.aliyun.com/alpine/v3.11/community \
    --repository http://mirrors.aliyun.com/alpine/v3.11/main

# 镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY="https://goproxy.cn,direct" \
    TZ Asia/Shanghai

# 移动到工作目录：/app
WORKDIR /build/app

# 复制项目中的 go.mod 和 go.sum文件并下载依赖信息
COPY go.mod .
COPY go.sum .
RUN go mod download

# 将代码复制到容器中
COPY . .

# Build the Go app
RUN go build -ldflags="-s -w" -o /app/snake .

# 创建一个小镜像
# Final stage
FROM debian:stretch-slim
# FROM alpine

WORKDIR /app

# 从builder镜像中把 /app 拷贝到当前目录
COPY --from=builder /app/snake /app/snake
COPY conf /app/conf

# Expose port 8080 to the outside world
EXPOSE 8080

# 需要运行的命令
ENTRYPOINT ["/app/snake"]
CMD ["-c", "conf/config.local.yaml"]
