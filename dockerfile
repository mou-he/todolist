# 第一阶段：构建（使用 Go 1.26）
FROM golang:1.26-alpine AS builder

ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制全部源码并编译
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -p 2 -ldflags="-s -w" -o server cmd/server/main.go

# 第二阶段：运行
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app

# 复制二进制文件
COPY --from=builder /app/server .

# 复制配置文件（宿主机 config/ 目录 → 容器 /app/config/）
COPY config/ ./config/

EXPOSE 8080
CMD ["./server"]