# ----------------------------------------------------
# 阶段 1: 构建环境
# ----------------------------------------------------
FROM --platform=$BUILDPLATFORM golang:1.23-bookworm AS builder

WORKDIR /app
ARG TARGETOS
ARG TARGETARCH

# 安装交叉编译器
RUN apt-get update && apt-get install -y \
    gcc-aarch64-linux-gnu \
    libc6-dev-arm64-cross \
    gcc-x86-64-linux-gnu \
    libc6-dev-amd64-cross \
    && rm -rf /var/lib/apt/lists/*

# 复制源码并处理依赖
COPY . .
RUN rm -f go.mod go.sum && \
    go mod init project-4869 && \
    go mod tidy

# 执行静态编译 (关键：静态链接 SQLite 库)
RUN if [ "$TARGETARCH" = "arm64" ]; then \
        export CC=aarch64-linux-gnu-gcc; \
    else \
        export CC=x86_64-linux-gnu-gcc; \
    fi && \
    CGO_ENABLED=1 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -v -ldflags="-s -w -extldflags '-static'" -o project4869 .

# ----------------------------------------------------
# 阶段 2: 运行环境
# ----------------------------------------------------
FROM mcr.microsoft.com/playwright:v1.41.0-jammy

WORKDIR /app
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY --from=builder /app/project4869 .
COPY static ./static

RUN mkdir -p data logs
ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright
EXPOSE 4869

CMD ["./project4869"]