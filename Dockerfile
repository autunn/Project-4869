# ----------------------------------------------------
# 阶段 1: 构建环境 (Builder)
# 使用标准 Go 镜像 (基于 Debian)，自带完整的 gcc 和 glibc
# 彻底解决 sqlite3 的 'pread64' 编译报错
# ----------------------------------------------------
FROM golang:1.21 AS builder

WORKDIR /app

# 标准镜像自带 git 和 gcc，无需手动安装

# 1. 复制源码
COPY . .

# 2. 自动修复 go.mod (防止手动上传时写错)
RUN rm -f go.mod go.sum
RUN go mod init project-4869
RUN go mod tidy

# 3. 编译
# CGO_ENABLED=1 开启 C 语言支持 (sqlite 需要)
# GOOS=linux 编译为 Linux 程序
RUN CGO_ENABLED=1 GOOS=linux go build -a -o project4869 .

# ----------------------------------------------------
# 阶段 2: 运行环境 (Runtime)
# 使用 Playwright 官方镜像 (基于 Ubuntu Jammy)
# ----------------------------------------------------
FROM mcr.microsoft.com/playwright:v1.41.0-jammy

WORKDIR /app

# 设置时区为上海
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 从第一阶段复制编译好的程序
COPY --from=builder /app/project4869 .
# 复制静态资源 (前端页面)
COPY static ./static

# 创建数据和日志目录
RUN mkdir -p data logs

# 设置 Playwright 浏览器路径
ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright

# 暴露端口
EXPOSE 4869

# 启动程序
CMD ["./project4869"]