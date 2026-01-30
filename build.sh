#!/usr/bin/env bash

# 编译脚本 - 用于编译前端、stat 和 xray-admin
# 编译后的文件会放到 dist 目录
# 前端编译产物会复制到 gin-vue-admin/server/resource/page 目录

set -e

Green="\033[32m"
Red="\033[31m"
Blue="\033[36m"
Yellow="\033[33m"
Font="\033[0m"

# 默认目标平台
TARGET_OS="linux"
TARGET_ARCH="amd64"
EXE_SUFFIX=""

# 显示帮助
show_help() {
  echo -e "${Blue}用法: ./build.sh [选项]${Font}"
  echo ""
  echo "选项:"
  echo "  -t, --target <platform>  目标平台: linux (默认), windows"
  echo "  -h, --help               显示帮助"
  echo ""
  echo "示例:"
  echo "  ./build.sh               编译 Linux 版本"
  echo "  ./build.sh -t linux      编译 Linux 版本"
  echo "  ./build.sh -t windows    编译 Windows 版本"
}

# 解析参数
while [[ $# -gt 0 ]]; do
  case $1 in
    -t|--target)
      case $2 in
        linux)
          TARGET_OS="linux"
          EXE_SUFFIX=""
          ;;
        windows)
          TARGET_OS="windows"
          EXE_SUFFIX=".exe"
          ;;
        *)
          echo -e "${Red}错误: 不支持的平台 '$2'，支持: linux, windows${Font}"
          exit 1
          ;;
      esac
      shift 2
      ;;
    -h|--help)
      show_help
      exit 0
      ;;
    *)
      echo -e "${Red}错误: 未知选项 '$1'${Font}"
      show_help
      exit 1
      ;;
  esac
done

echo -e "${Blue}========================================${Font}"
echo -e "${Blue}    编译前端、stat 和 xray-admin       ${Font}"
echo -e "${Blue}    目标平台: ${TARGET_OS}/${TARGET_ARCH}${Font}"
echo -e "${Blue}========================================${Font}"

# 创建输出目录
mkdir -p dist

# 编译前端
echo -e "\n${Green}[1/3] 编译前端...${Font}"
cd gin-vue-admin/web
if [ ! -d "node_modules" ]; then
  echo -e "${Yellow}安装前端依赖...${Font}"
  npm install --legacy-peer-deps
fi
npm run build
# 复制到 server/resource/page 目录 (后端会嵌入这些静态文件)
rm -rf ../server/resource/page
mkdir -p ../server/resource/page
cp -r dist/* ../server/resource/page/
cd ../..
echo -e "${Green}[OK] 前端编译完成${Font}"

# 编译 stat
echo -e "\n${Green}[2/3] 编译 stat...${Font}"
cd stat
go mod tidy
CGO_ENABLED=0 GOOS=${TARGET_OS} GOARCH=${TARGET_ARCH} go build -ldflags="-s -w" -o ../dist/stat${EXE_SUFFIX} .
cd ..
echo -e "${Green}[OK] stat 编译完成${Font}"

# 编译 xray-admin (gin-vue-admin server)
echo -e "\n${Green}[3/3] 编译 xray-admin...${Font}"
cd gin-vue-admin/server
go mod tidy
CGO_ENABLED=0 GOOS=${TARGET_OS} GOARCH=${TARGET_ARCH} go build -ldflags="-s -w" -o ../../dist/xray-admin${EXE_SUFFIX} .
cd ../..
echo -e "${Green}[OK] xray-admin 编译完成${Font}"

# 显示结果
echo -e "\n${Blue}========================================${Font}"
echo -e "${Green}编译完成！输出文件:${Font}"
ls -lh dist/
echo -e "${Blue}========================================${Font}"

if [ "$TARGET_OS" == "linux" ]; then
  echo -e "\n${Blue}下一步:${Font}"
  echo "1. 在 GitHub 创建 Release (如 v1.0.0)"
  echo "2. 上传 dist/stat 和 dist/xray-admin 到 Release"
else
  echo -e "\n${Blue}Windows 本地运行:${Font}"
  echo "1. 复制 dist/xray-admin.exe 到目标目录"
  echo "2. 配置 config.yaml"
  echo "3. 运行 xray-admin.exe"
fi
