#!/bin/bash

# 微博群消息发送工具 - 全平台打包脚本

set -e

echo "=== 开始构建并打包 ==="

# 清理旧文件
rm -rf dist release
mkdir -p dist

# 构建 Windows 版本
echo "构建 Windows 版本..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o dist/weibo-sender-windows.exe

# 构建 macOS 版本
echo "构建 macOS Intel 版本..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o dist/weibo-sender-macos-amd64

echo "构建 macOS Apple Silicon 版本..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -trimpath -o dist/weibo-sender-macos-arm64

# 构建 Linux 版本
echo "构建 Linux 版本..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o dist/weibo-sender-linux

# 复制必要文件
echo "复制文档文件..."
cp README.md dist/
cp QUICKSTART.md dist/
cp USAGE_EXAMPLES.md dist/

# 创建配置文件示例
echo "创建配置文件示例..."
cat > dist/config.example.json << 'EOF'
{
  "cookies": {
    "SCF": "",
    "SUB": "",
    "SUBP": "",
    "ALF": "",
    "_s_tentry": "weibo.com",
    "Apache": "",
    "SINAGLOBAL": "",
    "ULV": ""
  },
  "source": "209678993",
  "send_delay": 2
}
EOF

# 打包
echo "打包..."
cd dist
zip -r ../weibo-sender-all-platforms.zip .
cd ..

echo ""
echo "✓ 构建完成！"
echo ""
echo "输出文件: weibo-sender-all-platforms.zip"
ls -lh weibo-sender-all-platforms.zip
echo ""
echo "包含以下平台："
echo "  - Windows (amd64)"
echo "  - macOS Intel (amd64)"
echo "  - macOS Apple Silicon (arm64)"
echo "  - Linux (amd64)"

