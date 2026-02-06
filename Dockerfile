# ----------------------------------------------------
# 阶段 1: 构建阶段 (Builder)
# ----------------------------------------------------
FROM --platform=$BUILDPLATFORM golang:1.23-bookworm AS builder

WORKDIR /app

# 接收架构参数
ARG TARGETOS
ARG TARGETARCH

# 安装多架构 C 编译器 (关键：解决 SQLite 交叉编译报错)
RUN apt-get update && apt-get install -y \
    gcc-aarch64-linux-gnu \
    libc6-dev-arm64-cross \
    gcc-x86-64-linux-gnu \
    libc6-dev-amd64-cross \
    && rm -rf /var/lib/apt/lists/*

# 1. 复制源码
COPY . .

# 2. 强制初始化 Go 环境并清理冗余依赖
RUN rm -f go.mod go.sum && \
    go mod init project-4869 && \
    go mod tidy

# 3. 针对不同架构切换编译器并执行静态编译
RUN if [ "$TARGETARCH" = "arm64" ]; then \
        export CC=aarch64-linux-gnu-gcc; \
    else \
        export CC=x86_64-linux-gnu-gcc; \
    fi && \
    CGO_ENABLED=1 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -v -ldflags="-s -w -extldflags '-static'" -o project4869 .

# ----------------------------------------------------
# 阶段 2: 运行阶段 (Final Image)
# ----------------------------------------------------
# 使用 Playwright 镜像作为底座，确保支持网页抓取所需的浏览器环境
FROM mcr.microsoft.com/playwright:v1.41.0-jammy

WORKDIR /app

# 设置时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/project4869 .

# 复制前端静态资源
COPY static ./static

# 预创建数据和日志持久化目录
RUN mkdir -p data logs

# 设置 Playwright 浏览器路径
ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright

# 暴露 Web 端口
EXPOSE 4869

# 启动程序
CMD ["./project4869"]