#!/usr/bin/env bash

# 编译脚本 - 用于编译前端、stat 和 xray-admin
# 编译后的文件会放到 dist 目录

set -e

Green="\033[32m"
Red="\033[31m"
Blue="\033[36m"
Yellow="\033[33m"
Font="\033[0m"

echo -e "${Blue}========================================${Font}"
echo -e "${Blue}       Xray 项目编译脚本               ${Font}"
echo -e "${Blue}========================================${Font}"

# 选择目标平台
echo -e "\n${Green}选择目标平台:${Font}"
echo "1. Linux (默认)"
echo "2. Windows"
read -rp "请输入 (1/2): " platform_choice

case $platform_choice in
  2)
    TARGET_OS="windows"
    EXE_SUFFIX=".exe"
    ;;
  *)
    TARGET_OS="linux"
    EXE_SUFFIX=""
    ;;
esac

# 是否编译前端
echo -e "\n${Green}是否编译前端? (编译前端较耗时)${Font}"
read -rp "编译前端? (y/n, 默认n): " build_frontend

if [[ "$build_frontend" == "y" || "$build_frontend" == "Y" ]]; then
  SKIP_FRONTEND=false
else
  SKIP_FRONTEND=true
fi

echo -e "\n${Blue}========================================${Font}"
echo -e "${Blue}    目标平台: ${TARGET_OS}${Font}"
if [ "$SKIP_FRONTEND" = true ]; then
  echo -e "${Yellow}    跳过前端编译${Font}"
else
  echo -e "${Green}    编译前端${Font}"
fi
echo -e "${Blue}========================================${Font}"

# 创建输出目录
mkdir -p dist

# 编译前端
if [ "$SKIP_FRONTEND" = false ]; then
  echo -e "\n${Green}[1/3] 编译前端...${Font}"
  cd gin-vue-admin/web
  if [ ! -d "node_modules" ]; then
    echo -e "${Yellow}安装前端依赖...${Font}"
    npm install --legacy-peer-deps
  fi
  npm run build
  rm -rf ../server/resource/page
  mkdir -p ../server/resource/page
  cp -r dist/* ../server/resource/page/
  cd ../..
  echo -e "${Green}[OK] 前端编译完成${Font}"
else
  echo -e "\n${Yellow}[1/3] 跳过前端编译${Font}"
fi

# 编译 stat
echo -e "\n${Green}[2/3] 编译 stat...${Font}"
cd stat
go mod tidy
CGO_ENABLED=0 GOOS=${TARGET_OS} GOARCH=amd64 go build -ldflags="-s -w" -o ../dist/stat${EXE_SUFFIX} .
cd ..
echo -e "${Green}[OK] stat 编译完成${Font}"

# 编译 xray-admin (前端文件会通过 embed 嵌入到二进制中)
echo -e "\n${Green}[3/3] 编译 xray-admin...${Font}"
cd gin-vue-admin/server
go mod tidy
CGO_ENABLED=0 GOOS=${TARGET_OS} GOARCH=amd64 go build -ldflags="-s -w" -o ../../dist/xray-admin${EXE_SUFFIX} .
cd ../..
echo -e "${Green}[OK] xray-admin 编译完成 (前端已嵌入二进制)${Font}"

# 显示结果
echo -e "\n${Blue}========================================${Font}"
echo -e "${Green}编译完成！输出文件:${Font}"
ls -lh dist/
echo -e "${Blue}========================================${Font}"
