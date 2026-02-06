# ----------------------------------------------------
# 阶段 1: 构建环境 (Builder)
# 【关键修改】升级到 Go 1.23 (基于 Debian Bookworm)
# 彻底解决 'pread64' 报错 和 'requires go >= 1.23' 报错
# ----------------------------------------------------
FROM golang:1.23 AS builder

WORKDIR /app

# 1. 复制源码
COPY . .

# 2. 自动修复 go.mod
RUN rm -f go.mod go.sum
RUN go mod init project-4869
RUN go mod tidy

# 3. 编译
# CGO_ENABLED=1 开启 C 语言支持
# GOOS=linux 编译为 Linux 程序
RUN CGO_ENABLED=1 GOOS=linux go build -a -o project4869 .

# ----------------------------------------------------
# 阶段 2: 运行环境 (Runtime)
# 使用 Playwright 官方镜像 (基于 Ubuntu Jammy)
# ----------------------------------------------------
FROM mcr.microsoft.com/playwright:v1.41.0-jammy

WORKDIR /app

# 设置时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 从第一阶段复制编译好的程序
COPY --from=builder /app/project4869 .
COPY static ./static

# 创建目录
RUN mkdir -p data logs

# 环境变量
ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright

# 暴露端口
EXPOSE 4869

# 启动程序
CMD ["./project4869"]