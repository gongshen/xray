# Build Instructions - 代理服务器管理平台

## 概述
本文档提供项目各组件的构建指南。

---

## 1. 后端构建 (gin-vue-admin/server)

### 前置条件
- Go 1.18+ 已安装
- 配置好 GOPATH 和 GOPROXY

### 构建步骤

```bash
# 进入后端目录
cd gin-vue-admin/server

# 下载依赖
go mod tidy

# 编译
go build -o server.exe .

# 或 Linux/Mac
go build -o server .
```

### 验证
```bash
# 检查编译产物
ls -la server.exe  # Windows
ls -la server      # Linux/Mac
```

---

## 2. 前端构建 (gin-vue-admin/web)

### 前置条件
- Node.js 16+ 已安装
- npm 或 yarn 包管理器

### 构建步骤

```bash
# 进入前端目录
cd gin-vue-admin/web

# 安装依赖 (使用 legacy-peer-deps 解决依赖冲突)
npm install --legacy-peer-deps

# 生产构建
npm run build
```

### 验证
```bash
# 检查构建产物
ls -la dist/
```

### 构建产物
- `dist/` 目录包含静态资源
- 可部署到 Nginx 或其他 Web 服务器

---

## 3. Stat 程序构建 (stat)

### 前置条件
- Go 1.18+ 已安装
- 需要 xray-core 依赖

### 构建步骤

```bash
# 进入 stat 目录
cd stat

# 下载依赖
go mod tidy

# 编译
go build -o stat.exe .

# 或 Linux/Mac
go build -o stat .
```

### 验证
```bash
# 检查编译产物
ls -la stat.exe  # Windows
ls -la stat      # Linux/Mac
```

---

## 4. 一键构建脚本

项目根目录已有 `build.sh` 脚本，可用于一键构建：

```bash
# Linux/Mac
chmod +x build.sh
./build.sh

# Windows (Git Bash)
bash build.sh
```

---

## 构建状态检查清单

| 组件 | 命令 | 预期结果 |
|------|------|----------|
| 后端 | `go build` | 无错误，生成可执行文件 |
| 前端 | `npm run build` | 无错误，生成 dist/ 目录 |
| Stat | `go build` | 无错误，生成可执行文件 |

---

## 常见问题

### Q: Go 依赖下载失败
```bash
# 设置 GOPROXY
go env -w GOPROXY=https://goproxy.cn,direct
```

### Q: npm 依赖冲突
```bash
# 使用 legacy-peer-deps
npm install --legacy-peer-deps
```

### Q: 前端构建内存不足
```bash
# 增加 Node.js 内存限制
export NODE_OPTIONS="--max-old-space-size=4096"
npm run build
```
