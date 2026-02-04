#!/usr/bin/env bash

# Xray 项目编译脚本 - 用于编译前端、stat 和 xray-admin
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

# 第一层：是否编译前端
echo -e "\n${Green}[1/4] 是否编译前端?${Font}"
echo "前端代码有更新时必须选择编译，否则可以跳过以节省时间"
read -rp "编译前端? (y/n): " build_frontend

if [[ "$build_frontend" == "y" || "$build_frontend" == "Y" ]]; then
  BUILD_FRONTEND=true
else
  BUILD_FRONTEND=false
fi

# 第二层：是否深度清理缓存
echo -e "\n${Green}[2/4] 是否深度清理缓存?${Font}"
echo "深度清理会删除所有前端缓存"
echo "推荐在遇到缓存问题或首次编译时使用"
read -rp "深度清理缓存? (y/n): " deep_clean

if [[ "$deep_clean" == "y" || "$deep_clean" == "Y" ]]; then
  DEEP_CLEAN=true
else
  DEEP_CLEAN=false
fi

# 第三层：是否编译xray-admin
echo -e "\n${Green}[3/4] 是否编译xray-admin?${Font}"
echo "Web管理界面后端服务"
read -rp "编译xray-admin? (y/n): " build_admin

if [[ "$build_admin" == "y" || "$build_admin" == "Y" ]]; then
  BUILD_ADMIN=true
else
  BUILD_ADMIN=false
fi

# 第四层：是否编译stat
echo -e "\n${Green}[4/4] 是否编译stat?${Font}"
echo "流量统计服务"
read -rp "编译stat? (y/n): " build_stat

if [[ "$build_stat" == "y" || "$build_stat" == "Y" ]]; then
  BUILD_STAT=true
else
  BUILD_STAT=false
fi

# 显示编译配置
echo -e "\n${Blue}========================================${Font}"
echo -e "${Blue}           编译配置确认               ${Font}"
echo -e "${Blue}========================================${Font}"
echo -e "${Yellow}目标平台:     ${TARGET_OS}${Font}"
echo -e "${Yellow}编译前端:     $([ "$BUILD_FRONTEND" = true ] && echo "是" || echo "否")${Font}"
echo -e "${Yellow}深度清理:     $([ "$DEEP_CLEAN" = true ] && echo "是" || echo "否")${Font}"
echo -e "${Yellow}编译xray-admin: $([ "$BUILD_ADMIN" = true ] && echo "是" || echo "否")${Font}"
echo -e "${Yellow}编译stat:     $([ "$BUILD_STAT" = true ] && echo "是" || echo "否")${Font}"
echo -e "${Blue}========================================${Font}"

# 创建输出目录
mkdir -p dist

STEP=1
TOTAL_STEPS=$((${BUILD_FRONTEND} + ${DEEP_CLEAN} + ${BUILD_ADMIN} + ${BUILD_STAT}))

# 深度清理缓存
if [ "$DEEP_CLEAN" = true ]; then
  echo -e "\n${Green}[${STEP}/${TOTAL_STEPS}] 深度清理缓存...${Font}"
  
  # 清理embed目录
  echo -e "${Yellow}清理embed目录...${Font}"
  rm -rf gin-vue-admin/server/resource/page/*
  
  # 清理前端缓存和依赖
  echo -e "${Yellow}清理前端缓存和依赖...${Font}"
  cd gin-vue-admin/web
  rm -rf dist/
  rm -rf node_modules/
  rm -rf node_modules/.cache/
  rm -rf .vite/
  rm -rf .nuxt/
  rm -rf .output/
  rm -rf .temp/
  rm -rf package-lock.json
  
  # 清理更多前端缓存
  echo -e "${Yellow}清理更多前端缓存...${Font}"
  rm -rf .cache/
  rm -rf .parcel-cache/
  rm -rf .next/
  rm -rf build/
  rm -rf coverage/
  rm -rf .nyc_output/
  rm -rf .eslintcache
  rm -rf .stylelintcache
  
  # 清理npm缓存
  echo -e "${Yellow}清理npm缓存...${Font}"
  npm cache clean --force 2>/dev/null || true
  
  cd ../..
  
  # 清理Go构建缓存
  echo -e "${Yellow}清理Go构建缓存...${Font}"
  go clean -cache -modcache -i -r 2>/dev/null || true
  
  # 清理dist目录
  echo -e "${Yellow}清理dist目录...${Font}"
  rm -rf dist/*
  
  echo -e "${Green}✅ 深度清理完成${Font}"
  STEP=$((STEP + 1))
fi

# 编译前端
if [ "$BUILD_FRONTEND" = true ]; then
  echo -e "\n${Green}[${STEP}/${TOTAL_STEPS}] 编译前端...${Font}"
  cd gin-vue-admin/web
  
  # 强制清理前端构建缓存
  echo -e "${Yellow}强制清理前端构建缓存...${Font}"
  rm -rf dist/
  rm -rf .vite/
  rm -rf .cache/
  rm -rf node_modules/.cache/
  
  # 安装依赖
  if [ ! -d "node_modules" ]; then
    echo -e "${Yellow}安装前端依赖...${Font}"
    npm install --legacy-peer-deps
    if [ $? -ne 0 ]; then
      echo -e "${Red}❌ 前端依赖安装失败${Font}"
      exit 1
    fi
  fi
  
  echo -e "${Yellow}构建前端项目...${Font}"
  # 强制重新构建，忽略缓存
  npm run build -- --force
  if [ $? -ne 0 ]; then
    echo -e "${Red}❌ 前端构建失败，尝试清理缓存后重新构建...${Font}"
    rm -rf dist/ .vite/ .cache/ node_modules/.cache/
    npm run build
    if [ $? -ne 0 ]; then
      echo -e "${Red}❌ 前端构建失败${Font}"
      exit 1
    fi
  fi
  
  echo -e "${Yellow}复制前端文件到embed目录...${Font}"
  rm -rf ../server/resource/page
  mkdir -p ../server/resource/page
  cp -r dist/* ../server/resource/page/
  
  # 验证文件复制
  if [ ! -f "../server/resource/page/index.html" ]; then
    echo -e "${Red}❌ 前端文件复制失败${Font}"
    exit 1
  fi
  
  cd ../..
  echo -e "${Green}✅ 前端编译完成${Font}"
  STEP=$((STEP + 1))
fi

# 编译stat
if [ "$BUILD_STAT" = true ]; then
  echo -e "\n${Green}[${STEP}/${TOTAL_STEPS}] 编译 stat 服务...${Font}"
  cd stat
  
  echo -e "${Yellow}整理Go模块依赖...${Font}"
  go mod tidy
  if [ $? -ne 0 ]; then
    echo -e "${Red}❌ stat模块依赖整理失败${Font}"
    exit 1
  fi
  
  echo -e "${Yellow}编译stat二进制文件...${Font}"
  CGO_ENABLED=0 GOOS=${TARGET_OS} GOARCH=amd64 go build -ldflags="-s -w" -o ../dist/stat${EXE_SUFFIX} .
  if [ $? -ne 0 ]; then
    echo -e "${Red}❌ stat编译失败${Font}"
    exit 1
  fi
  
  cd ..
  echo -e "${Green}✅ stat 编译完成${Font}"
  STEP=$((STEP + 1))
fi

# 编译xray-admin
if [ "$BUILD_ADMIN" = true ]; then
  echo -e "\n${Green}[${STEP}/${TOTAL_STEPS}] 编译 xray-admin...${Font}"
  cd gin-vue-admin/server
  
  echo -e "${Yellow}整理Go模块依赖...${Font}"
  go mod tidy
  if [ $? -ne 0 ]; then
    echo -e "${Red}❌ xray-admin模块依赖整理失败${Font}"
    exit 1
  fi
  
  echo -e "${Yellow}编译xray-admin二进制文件 (前端已嵌入)...${Font}"
  CGO_ENABLED=0 GOOS=${TARGET_OS} GOARCH=amd64 go build -ldflags="-s -w" -o ../../dist/xray-admin${EXE_SUFFIX} .
  if [ $? -ne 0 ]; then
    echo -e "${Red}❌ xray-admin编译失败${Font}"
    exit 1
  fi
  
  cd ../..
  echo -e "${Green}✅ xray-admin 编译完成${Font}"
  STEP=$((STEP + 1))
fi

# 显示结果
echo -e "\n${Blue}========================================${Font}"
echo -e "${Green}           编译完成！                 ${Font}"
echo -e "${Blue}========================================${Font}"

if [ -d "dist" ] && [ "$(ls -A dist 2>/dev/null)" ]; then
  echo -e "${Green}输出文件:${Font}"
  ls -lh dist/
else
  echo -e "${Yellow}没有生成任何文件${Font}"
fi

# 使用提示
echo -e "\n${Blue}========================================${Font}"
echo -e "${Green}使用提示:${Font}"

if [ "$BUILD_ADMIN" = true ]; then
  echo -e "${Yellow}启动xray-admin:${Font}"
  echo -e "  cd dist"
  echo -e "  ./xray-admin${EXE_SUFFIX}"
fi

if [ "$BUILD_STAT" = true ]; then
  echo -e "${Yellow}启动stat服务:${Font}"
  echo -e "  cd dist"
  echo -e "  ./stat${EXE_SUFFIX}"
fi

if [ "$DEEP_CLEAN" = true ] || [ "$BUILD_FRONTEND" = true ]; then
  echo -e "\n${Yellow}重要提醒:${Font}"
  if [ "$DEEP_CLEAN" = true ]; then
    echo -e "  - 已深度清理缓存，请重启服务并清除浏览器缓存"
    echo -e "  - 建议按 Ctrl+F5 或 Ctrl+Shift+R 强制刷新浏览器"
    echo -e "  - 或者在浏览器开发者工具中右键刷新按钮选择'清空缓存并硬性重新加载'"
  fi
  if [ "$BUILD_FRONTEND" = true ]; then
    echo -e "  - 前端文件已嵌入到xray-admin二进制文件中"
    echo -e "  - 如需更新前端，必须重新编译"
    echo -e "  - 请重启xray-admin服务以加载新的前端文件"
  fi
fi

echo -e "${Blue}========================================${Font}"
