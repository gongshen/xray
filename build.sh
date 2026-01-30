#!/usr/bin/env bash

# 编译脚本 - 用于编译 stat 和 xray-admin (Linux amd64)
# 编译后的文件会放到 dist 目录

set -e

Green="\033[32m"
Red="\033[31m"
Blue="\033[36m"
Font="\033[0m"

echo -e "${Blue}========================================${Font}"
echo -e "${Blue}       编译 stat 和 xray-admin         ${Font}"
echo -e "${Blue}========================================${Font}"

# 创建输出目录
mkdir -p dist

# 编译 stat
echo -e "\n${Green}[1/2] 编译 stat...${Font}"
cd stat
go mod tidy
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../dist/stat .
cd ..
echo -e "${Green}[OK] stat 编译完成${Font}"

# 编译 xray-admin (gin-vue-admin server)
echo -e "\n${Green}[2/2] 编译 xray-admin...${Font}"
cd gin-vue-admin/server
go mod tidy
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../../dist/xray-admin .
cd ../..
echo -e "${Green}[OK] xray-admin 编译完成${Font}"

# 显示结果
echo -e "\n${Blue}========================================${Font}"
echo -e "${Green}编译完成！输出文件:${Font}"
ls -lh dist/
echo -e "${Blue}========================================${Font}"

echo -e "\n${Blue}下一步:${Font}"
echo "1. 在 GitHub 创建 Release (如 v1.0.0)"
echo "2. 上传 dist/stat 和 dist/xray-admin 到 Release"
echo "3. 上传 deploy/install.sh 到仓库"
