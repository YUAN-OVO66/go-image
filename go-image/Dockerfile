# 构建阶段
FROM golang:1.20-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的构建工具
RUN apk add --no-cache git

# 复制go.mod和go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o go-image ./cmd/server

# 运行阶段
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 设置时区为Asia/Shanghai
ENV TZ=Asia/Shanghai

# 创建非root用户
RUN adduser -D -g '' appuser

# 创建必要的目录
RUN mkdir -p /app/static/uploads/thumbnails && \
    chown -R appuser:appuser /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/go-image /app/
COPY --from=builder /app/static /app/static
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/configs /app/configs

# 切换到非root用户
USER appuser

# 设置工作目录
WORKDIR /app

# 暴露端口
EXPOSE 28080

# 启动应用
CMD ["./go-image"]