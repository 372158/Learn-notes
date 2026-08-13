# 第一阶段：构建阶段
FROM golang:1.23-alpine AS builder

# 设置工作目录
WORKDIR /app

# 先复制依赖文件，利用 Docker 缓存加速构建
COPY go.mod go.sum ./
RUN go mod download

# 复制所有源码并编译
COPY . .
RUN go build -o shortlink cmd/server/main.go

# 第二阶段：运行阶段
FROM alpine:latest

WORKDIR /app

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/shortlink .

# 暴露端口（与 config.yaml 中的 port 对应）
EXPOSE 8080

# 运行程序
CMD ["./shortlink"]