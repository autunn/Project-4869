# ----------------------------------------------------
# 阶段 1: 构建 Go 二进制文件
# ----------------------------------------------------
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装必要的系统库 (CGO 需要 gcc)
RUN apk add --no-cache gcc musl-dev

# 1. 【重要修改】先复制所有源代码
COPY . .

# 2. 【重要修改】现在代码都在了，再让 Go 自动整理依赖
# 这会自动生成 go.sum 并下载缺少的包
RUN go mod tidy
RUN go mod download

# 3. 编译静态二进制文件
# CGO_ENABLED=1 是必须的，因为我们要用 sqlite
RUN CGO_ENABLED=1 GOOS=linux go build -a -o project4869 .

# ----------------------------------------------------
# 阶段 2: 运行时环境 (包含 Playwright 浏览器)
# ----------------------------------------------------
FROM mcr.microsoft.com/playwright:v1.41.0-jammy

WORKDIR /app

# 设置国内时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 从 builder 阶段复制编译好的程序
COPY --from=builder /app/project4869 .
COPY static ./static

# 创建必要的目录
RUN mkdir -p data logs

# 环境变量
ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright

# 暴露端口
EXPOSE 4869

# 启动命令
CMD ["./project4869"]