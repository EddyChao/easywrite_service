# 第一阶段：编译阶段
FROM golang:1.22-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制其余的源代码
COPY . .
RUN go mod download
RUN go mod tidy
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12
RUN swag init --parseDependency --parseInternal
RUN go build -o myapp .

# 第二阶段：运行阶段
FROM alpine:latest

# 设置工作目录
WORKDIR /root/

RUN apk add --no-cache tzdata

# 从编译阶段复制二进制文件
COPY --from=builder /app/myapp .

# 运行二进制文件
CMD ["./myapp","-c","/etc/config/service_config.json5"]
