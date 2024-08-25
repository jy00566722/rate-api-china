# 使用官方的 Golang 镜像用作构建环境
FROM golang:1.19 AS builder
ARG GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 if available go.sum.
COPY go.mod ./
COPY go.sum ./

# 下载所有依赖项
RUN go mod download

# 复制源代码到容器中
COPY . .

# 构建程序，关闭 CGO 并进行静态链接
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp-rate .

# 使用 scratch 创建一个干净的最小容器
FROM alpine:3.18.2

# 从构建器镶嵌阶段复制可执行文件
COPY --from=builder /app/myapp-rate /

# 运行程序时监听的端口
EXPOSE 8001

# 运行程序
CMD ["/myapp-rate"]