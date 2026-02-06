# ----------------------------------------------------
# 阶段 1: 准备工具链
# ----------------------------------------------------
FROM --platform=$BUILDPLATFORM tonistiigi/xx AS xx
FROM --platform=$BUILDPLATFORM golang:1.23 AS builder

# 复制交叉编译助手工具
COPY --from=xx / /

WORKDIR /app

# 获取构建的目标架构
ARG TARGETPLATFORM

# 安装编译 SQLite 所需的 C 开发库
RUN apt-get update && apt-get install -y binutils gcc g++

# 自动安装并配置对应架构的编译器
RUN xx-apt-get install -y gcc libc6-dev

# 复制源码
COPY . .

# 初始化并下载依赖
RUN rm -f go.mod go.sum && \
    go mod init project-4869 && \
    go mod tidy

# 使用 xx-go 进行编译
# xx-go 会自动处理 CGO_ENABLED=1, GOOS, GOARCH 和对应的 CC 编译器
RUN xx-go build -ldflags="-s -w" -o project4869 . && \
    xx-verify project4869

# ----------------------------------------------------
# 阶段 2: 运行环境 (Playwright)
# ----------------------------------------------------
FROM mcr.microsoft.com/playwright:v1.41.0-jammy

WORKDIR /app

# 设置时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 复制二进制文件和前端文件
COPY --from=builder /app/project4869 .
COPY static ./static

# 创建必要目录
RUN mkdir -p data logs

# 环境配置
ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright
EXPOSE 4869

CMD ["./project4869"]