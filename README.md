# Xray 管理系统

基于 [gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin) 开发的 Xray 代理服务管理系统。

## 项目结构

```
├── gin-vue-admin/       # 管理后台 (前端 + 后端)
│   ├── server/          # Go 后端服务 (xray-admin)
│   └── web/             # Vue 前端
├── stat/                # 流量统计服务
├── deploy/              # 部署脚本
│   └── install.sh       # 一键安装脚本
├── build.sh             # 编译脚本
└── dist/                # 编译输出目录
```

## 功能特性

- **服务器管理**: 添加、编辑、删除代理服务器
- **用户绑定管理**: 管理用户与服务器的绑定关系
- **流量统计**: 实时统计用户流量使用情况
- **系统信息**: 收集服务器系统信息 (CPU、内存、磁盘等)
- **每服务器独立端口**: 支持为每台服务器配置独立的统计服务端口

## 系统要求

- CentOS 7+ / Debian 9+ / Ubuntu 18+
- MySQL 8.0+
- Go 1.19+ (仅编译需要)
- Node.js 16+ (仅前端开发需要)

## 快速安装

### 代理服务器端 (安装 Xray + Stat)

```bash
# 下载安装脚本
wget -O install.sh https://raw.githubusercontent.com/gongshen/xray/main/deploy/install.sh
chmod +x install.sh
bash install.sh
```

安装菜单选项：
1. 安装 Xray (含依赖和优化)
2. 配置 Xray
3. 安装 BBR 加速脚本
4. 安装 Stat 服务
5. 安装 xray-admin 管理端
6. 安装 MySQL
7. 初始化管理端数据库
8. 配置 iptables 防火墙
9. 安装 acme.sh (SSL证书工具)
10. 续期 SSL 证书

### 管理端 (安装 xray-admin)

在管理服务器上运行安装脚本，选择菜单 5、6、7 完成管理端安装。

## 手动编译

```bash
# 编译 stat 和 xray-admin
./build.sh

# 编译后的文件在 dist/ 目录
ls dist/
# stat        - 流量统计服务
# xray-admin  - 管理后台服务
```

## 配置说明

### Stat 服务

```bash
# 启动参数
stat -port 56611 -level info

# 参数说明
-port   监听端口 (默认: 56611)
-level  日志级别 (默认: info)

# 环境变量
REMOTE_IP  允许访问的管理端 IP 地址
```

### xray-admin 服务

配置文件: `/usr/local/etc/xray-admin/config.yaml`

```yaml
system:
  addr: 8888              # 管理端口

mysql:
  path: "127.0.0.1"       # 数据库地址
  port: "3306"            # 数据库端口
  db-name: "gva"          # 数据库名
  username: "root"        # 用户名
  password: ""            # 密码

stat_port: 56611          # 默认统计服务端口
```

## 服务管理

```bash
# Stat 服务
systemctl start stat
systemctl stop stat
systemctl restart stat
systemctl status stat

# xray-admin 服务
systemctl start xray_admin
systemctl stop xray_admin
systemctl restart xray_admin
systemctl status xray_admin

# Xray 服务
systemctl start xray
systemctl stop xray
systemctl restart xray
systemctl status xray
```

## 默认端口

| 服务 | 端口 | 说明 |
|------|------|------|
| xray-admin | 8888 | 管理后台 Web 端口 |
| stat | 56611 | 流量统计服务端口 (可配置) |
| xray | 80 | 代理服务端口 |
| xray api | 11111 | Xray 内部 API 端口 |

## 开发

### 后端开发

```bash
cd gin-vue-admin/server
go mod tidy
go run main.go
```

### 前端开发

```bash
cd gin-vue-admin/web
npm install --legacy-peer-deps
npm run serve
```

### 前端构建

```bash
cd gin-vue-admin/web
npm run build
```

## License

MIT
