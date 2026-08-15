# Dockerfile
# 多阶段构建 Go 应用镜像
# 学习要点：多阶段构建能大幅减小镜像体积
#   - 第一阶段（builder）：用 Go 镜像编译二进制文件
#   - 第二阶段：用精简镜像运行二进制文件，不含编译工具
#
# 为什么用多阶段构建：
# 1. Go 镜像约 800MB，alpine 镜像约 10MB
# 2. 运行时不需要 Go 编译器和源码
# 3. 最终镜像小，拉取快，安全风险低

# ===== 第一阶段：构建阶段 =====
FROM golang:1.21-alpine AS builder

# 设置工作目录
# 学习要点：WORKDIR 创建并切换到指定目录，后续命令都在这执行
WORKDIR /app

# 先复制依赖文件（利用 Docker 缓存层）
# 学习要点：先复制 go.mod/go.sum 再下载依赖，利用缓存加速构建
# 只要依赖不变，下载依赖这一层会缓存，不用每次重新下载
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译二进制文件
# 学习要点：
# - CGO_ENABLED=0: 禁用 CGO，编译出纯静态二进制（可在 alpine 运行）
# - GOOS=linux: 目标系统 Linux（容器是 Linux）
# - -o app: 输出文件名 app
# - -ldflags="-s -w": 去掉调试信息，减小二进制体积
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app .


# ===== 第二阶段：运行阶段 =====
FROM alpine:latest

# 安装 ca-certificates（HTTPS 请求需要证书）
# 学习要点：alpine 镜像很精简，很多基础组件没有，按需安装
RUN apk --no-cache add ca-certificates

WORKDIR /app

# 从 builder 阶段复制编译好的二进制文件
# 学习要点：--from=builder 指定从哪个阶段复制
COPY --from=builder /app/app .
# 复制配置文件
COPY --from=builder /app/config ./config

# 暴露端口（仅声明，实际映射在 docker-compose 或 docker run -p 中配置）
EXPOSE 8080

# 启动命令
# 学习要点：CMD 容器启动时执行的命令
CMD ["./app"]
