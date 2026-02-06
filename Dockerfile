# ----------------------------------------------------
# 阶段 1: 构建环境 (Builder)
# 关键点：使用 golang:1.23 (基于 Debian)
# 1. 解决 'go >= 1.23' 报错
# 2. 解决 ARM 架构下 sqlite3 'pread64' 报错
# ----------------------------------------------------
FROM golang:1.23 AS builder

WORKDIR /app

# 1. 复制所有源代码
COPY . .

# 2. 强制重新生成 go.mod
# 防止之前的配置有误
RUN rm -f go.mod go.sum
RUN go mod init project-4869
RUN go mod tidy

# 3. 编译
# 这里的 GOARCH 会由 Docker Buildx 自动注入（构建 ARM 时自动为 arm64）
# 所以我们只需要指定 OS 和 CGO 即可
RUN CGO_ENABLED=1 GOOS=linux go build -a -o project4869 .

# ----------------------------------------------------
# 阶段 2: 运行环境 (Runtime)
# ----------------------------------------------------
FROM mcr.microsoft.com/playwright:v1.41.0-jammy

WORKDIR /app

# 设置时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 复制编译产物
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