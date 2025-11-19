# 多阶段构建 Dockerfile
# 第一阶段：构建阶段
FROM golang:1.25.4-bookworm AS builder

# 安装必要的构建工具（Debian 使用 apt）
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    git \
    make \
    gcc \
    libc6-dev && \
    rm -rf /var/lib/apt/lists/*

# 设置工作目录
WORKDIR /build

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o team-navigation main.go

# 第二阶段：运行阶段
FROM alpine:latest

# 安装必要的运行时依赖（SQLite 需要，wget 用于健康检查）
RUN apk --no-cache add ca-certificates sqlite wget

# 创建应用用户
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /build/team-navigation .

# 复制源代码和资源文件（应用运行时需要）
COPY --from=builder /build/src ./src
COPY --from=builder /build/pages ./pages
COPY --from=builder /build/http.zy .
COPY --from=builder /build/database ./database

# 创建数据库目录并设置权限
RUN mkdir -p /app/data && \
    chown -R appuser:appuser /app

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8080

# 设置环境变量
ENV GIN_MODE=release

# 启动应用
CMD ["./team-navigation"]

