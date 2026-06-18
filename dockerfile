# 第一阶段：构建（使用 Alpine 版本，镜像更小）
FROM golang:1.21-alpine AS builder

# 设置国内代理（可选，加速下载）
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app

# 复制依赖文件并下载（利用 Docker 缓存）
COPY go.mod go.sum ./
RUN go mod download

# 复制全部源码并编译（关闭 CGO 以生成静态二进制文件）
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

# 第二阶段：运行（使用极简的 Alpine 基础镜像）
FROM alpine:latest

# 安装 ca-certificates 以便访问 HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# 从构建阶段复制二进制文件和配置文件目录
COPY --from=builder /app/server .
COPY --from=builder /app/configs ./configs

# 暴露端口（与 config.yaml 保持一致）
EXPOSE 8080

# 启动服务
CMD ["./server"]