# ----------------------------------------------------
# 阶段 1: 构建环境
# ----------------------------------------------------
FROM --platform=$BUILDPLATFORM golang:1.23 AS builder

WORKDIR /app

# 获取构建的目标架构（由 Buildx 自动注入）
ARG TARGETOS
ARG TARGETARCH

# 安装编译所需的 C 依赖
# Debian 基础镜像稳定性最高
RUN apt-get update && apt-get install -y gcc-aarch64-linux-gnu gcc-x86-64-linux-gnu g++-aarch64-linux-gnu g++-x86-64-linux-gnu libc6-dev-arm64-cross

# 复制源码
COPY . .

# 初始化并下载依赖
RUN rm -f go.mod go.sum && \
    go mod init project-4869 && \
    go mod tidy

# 针对不同架构设置不同的 C 编译器 (关键步骤)
RUN if [ "$TARGETARCH" = "arm64" ]; then \
        export CC=aarch64-linux-gnu-gcc; \
    else \
        export CC=x86_64-linux-gnu-gcc; \
    fi && \
    CGO_ENABLED=1 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="-s -w" -o project4869 .

# ----------------------------------------------------
# 阶段 2: 运行环境
# ----------------------------------------------------
FROM mcr.microsoft.com/playwright:v1.41.0-jammy

WORKDIR /app

# 设置时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 复制产物
COPY --from=builder /app/project4869 .
COPY static ./static

# 数据与日志目录
RUN mkdir -p data logs

# 设置浏览器路径
ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright

EXPOSE 4869

CMD ["./project4869"]