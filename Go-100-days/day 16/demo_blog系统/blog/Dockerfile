# 1. 使用官方 Go 镜像作为构建环境
FROM golang:1.21-alpine AS builder

# 2. 设置工作目录
WORKDIR /app

# 3. 复制 go.mod 和 go.sum，并下载依赖（利用 Docker 层缓存）
COPY go.mod go.sum ./
RUN go mod download

# 4. 复制所有源代码
COPY . .

# 5. 编译为可执行文件（静态编译，CGO_ENABLED=0）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o blog-server ./main.go

# 6. 使用极小的 Alpine 镜像作为运行环境
FROM alpine:latest

# 7. 安装 ca-certificates（支持 HTTPS）
RUN apk --no-cache add ca-certificates

# 8. 设置工作目录
WORKDIR /root/

# 9. 从构建阶段复制编译好的可执行文件
COPY --from=builder /app/blog-server .

# 10. 复制配置文件（如果需要）
COPY config.yaml .

# 11. 暴露端口（与 config.yaml 中的端口一致）
EXPOSE 8080

# 12. 启动命令
CMD ["./blog-server"]