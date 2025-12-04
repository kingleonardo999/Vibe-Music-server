# ---- Build Stage ----
# 使用官方的 Go 镜像作为构建环境
# 使用 alpine 版本以保持镜像体积较小
FROM golang:1.24-alpine AS builder

# 设置必要的环境变量
# CGO_ENABLED=0 禁用 CGO，以构建静态链接的二进制文件
# GOOS=linux 指定目标操作系统
ENV CGO_ENABLED=0 GOOS=linux GO111MODULE=on

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
# 这样做可以利用 Docker 的层缓存机制，只有在依赖变更时才重新下载
COPY ../go.mod ../go.sum ./
RUN go mod download

# 复制所有源代码到工作目录
COPY ../. .

# 编译 Go 应用
# -o 指定输出文件名
# -ldflags "-s -w" 去除调试信息，减小二进制文件大小
RUN go build -o vibe-music-server -ldflags="-s -w" ./cmd/server/main.go


# ---- Runtime Stage ----
# 使用一个非常小的基础镜像
FROM alpine:latest

# 安装 ca-certificates 以支持 HTTPS 请求，以及 tzdata 设置时区
RUN apk --no-cache add ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/vibe-music-server .

# 暴露应用程序运行的端口
# 默认端口为 8080，可以通过环境变量 SERVER_PORT 覆盖
EXPOSE 8080

# 设置容器启动时执行的命令
ENTRYPOINT ["./vibe-music-server"]

