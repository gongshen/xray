# Xray 管理系统

基于 [gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin) 开发的现代化 Xray 代理服务管理系统，提供完整的代理服务器管理、用户绑定、流量统计和系统监控功能。

## ✨ 功能特性

### 🖥️ 管理后台功能
- **服务器管理**: 添加、编辑、删除代理服务器，支持独立统计端口配置
- **用户绑定管理**: 管理用户与服务器的绑定关系，支持限流控制
- **配置分享**: 生成 Shadowrocket、V2rayN 等客户端配置二维码和链接
- **代理重启**: 远程重启代理服务，无需登录服务器

### 📊 流量统计功能
- **实时流量监控**: 实时统计用户流量使用情况
- **流量趋势图**: 可视化展示流量使用趋势
- **流量排行榜**: TOP 10 用户流量使用排行
- **详细流量记录**: 按时间、用户、服务器查看详细流量数据
- **多维度统计**: 支持按日期范围、服务器IP、用户等维度统计

### 🔐 权限管理
- **管理员权限**: 查看所有用户数据和服务器管理
- **用户权限**: 查看个人流量数据和配置分享
- **角色分离**: 代理管理和个人代理功能分离

### 📱 响应式设计
- **移动端适配**: 完美支持手机、平板等移动设备
- **现代化UI**: 基于 Element Plus 的现代化界面设计
- **暗色主题**: 支持暗色模式切换

## 🏗️ 项目结构

```
├── gin-vue-admin/           # 管理后台
│   ├── server/              # Go 后端服务 (xray-admin)
│   │   ├── api/             # API 接口层
│   │   ├── model/           # 数据模型
│   │   ├── service/         # 业务逻辑层
│   │   ├── router/          # 路由配置
│   │   └── config.yaml      # 配置文件
│   └── web/                 # Vue 前端
│       ├── src/view/v2ray_admin/  # 管理员功能页面
│       ├── src/view/v2ray/        # 用户功能页面
│       └── src/api/               # API 接口定义
├── stat/                    # 流量统计服务
│   ├── business/            # 业务逻辑
│   ├── server/              # HTTP 服务器
│   └── main.go              # 入口文件
├── deploy/                  # 部署脚本
│   └── install.sh           # 一键安装脚本
├── build.sh                 # 编译脚本
└── dist/                    # 编译输出目录
```

## 🚀 快速开始

### 系统要求

- **操作系统**: CentOS 7+ / Debian 9+ / Ubuntu 18+
- **数据库**: MySQL 8.0+
- **开发环境**: Go 1.19+ / Node.js 16+ (仅开发需要)

### 一键安装

#### 代理服务器端安装

```bash
# 下载安装脚本
wget -O install.sh https://raw.githubusercontent.com/gongshen/xray/main/deploy/install.sh
chmod +x install.sh
bash install.sh
```

**安装菜单选项：**
1. **安装 Xray** - 安装 Xray 核心及依赖，系统优化
2. **配置 Xray** - 配置 Xray 服务参数
3. **安装 BBR 加速** - 安装 BBR 网络加速脚本
4. **安装 Stat 服务** - 安装流量统计服务
5. **安装 xray-admin** - 安装管理后台服务
6. **安装 MySQL** - 安装 MySQL 数据库
7. **初始化数据库** - 初始化管理端数据库结构
8. **配置防火墙** - 配置 iptables 防火墙规则
9. **安装 SSL 工具** - 安装 acme.sh SSL 证书工具
10. **续期 SSL 证书** - 自动续期 SSL 证书
11. **配置 Xray 日志轮转** - 为 Xray access.log/error.log 配置按天轮转
12. **分析用户流量明细** - 按用户 tag/id 查询 stat SQLite 采集周期流量和 Xray access.log 访问明细，单次时间区间最多 1 天

#### 管理端安装

在管理服务器上运行安装脚本，依次选择菜单 **5 → 6 → 7** 完成管理端安装。

### 手动编译

```bash
# 克隆项目
git clone https://github.com/gongshen/xray.git
cd xray

# 使用交互式编译脚本
./build.sh

# 编译选项：
# 1. 是否编译前端? (y/n)
# 2. 是否深度清理缓存? (y/n) 
# 3. 是否编译xray-admin? (y/n)
# 4. 是否编译stat? (y/n)

# 编译后的文件在 dist/ 目录
ls dist/
# stat        - 流量统计服务
# xray-admin  - 管理后台服务
```

## ⚙️ 配置说明

### Stat 服务配置

```bash
# 启动参数
./stat -port 56611 -level info

# 参数说明
-port   监听端口 (默认: 56611)
-level  日志级别 (debug/info/warn/error，默认: info)

# 环境变量
export REMOTE_IP="管理端IP地址"  # 允许访问的管理端 IP
```

### xray-admin 服务配置

配置文件位置: `/usr/local/etc/xray-admin/config.yaml`

```yaml
# 系统配置
system:
  addr: 8888              # 管理端口
  db-type: mysql          # 数据库类型

# MySQL 数据库配置
mysql:
  path: "127.0.0.1"       # 数据库地址
  port: "3306"            # 数据库端口
  db-name: "gva"          # 数据库名
  username: "root"        # 用户名
  password: "your_password" # 密码
  max-idle-conns: 10      # 最大空闲连接数
  max-open-conns: 100     # 最大打开连接数

# JWT 配置
jwt:
  signing-key: "your_secret_key"
  expires-time: 7d        # 过期时间
  buffer-time: 1d         # 缓冲时间

# 日志配置
zap:
  level: info             # 日志级别
  prefix: '[xray-admin]'  # 日志前缀
  director: log           # 日志目录

# 节点统计配置
stat_port: 56611          # 节点 stat 默认端口
traffic_collect_interval: 1h # xray-admin 拉取节点流量的间隔，未配置默认 1h
sysinfo_collect_interval: 5m # xray-admin 刷新在线状态/CPU/内存/磁盘的间隔，未配置默认 5m；10 分钟未更新显示离线
```

## 🔧 服务管理

### 系统服务控制

```bash
# Stat 流量统计服务
systemctl start stat      # 启动
systemctl stop stat       # 停止
systemctl restart stat    # 重启
systemctl status stat     # 查看状态
systemctl enable stat     # 开机自启

# xray-admin 管理服务
systemctl start xray_admin
systemctl stop xray_admin
systemctl restart xray_admin
systemctl status xray_admin
systemctl enable xray_admin

# Xray 代理服务
systemctl start xray
systemctl stop xray
systemctl restart xray
systemctl status xray
systemctl enable xray
```

### 日志查看

```bash
# 查看服务日志
journalctl -u stat -f          # 实时查看 stat 日志
journalctl -u xray_admin -f     # 实时查看 xray-admin 日志
journalctl -u xray -f           # 实时查看 xray 日志

# 查看应用日志
tail -f /usr/local/etc/xray-admin/log/server.log  # xray-admin 应用日志
```

## 🌐 端口说明

| 服务 | 默认端口 | 说明 | 可配置 |
|------|----------|------|--------|
| xray-admin | 8888 | 管理后台 Web 端口 | ✅ |
| stat | 56611 | 流量统计服务端口 | ✅ |
| xray | 80 | 代理服务端口 | ✅ |
| xray-api | 11111 | Xray 内部 API 端口 | ✅ |
| mysql | 3306 | MySQL 数据库端口 | ✅ |

## 🎯 使用指南

### 管理员功能

1. **服务器管理**
   - 访问 `http://your-server:8888` 登录管理后台
   - 添加代理服务器，配置IP、端口、统计端口等
   - 设置流量配额和重置时间

2. **用户绑定管理**
   - 创建用户与服务器的绑定关系
   - 支持限流控制和解除限流
   - 生成配置分享链接和二维码

3. **流量统计分析**
   - 查看流量趋势图和排行榜
   - 按时间范围、服务器、用户筛选数据
   - 导出流量报表

### 用户功能

1. **个人流量查看**
   - 查看个人流量使用情况
   - 流量趋势分析

2. **配置获取**
   - 获取个人代理配置
   - 扫描二维码或复制链接导入客户端

## 🛠️ 开发指南

### 后端开发

```bash
# 进入后端目录
cd gin-vue-admin/server

# 安装依赖
go mod tidy

# 运行开发服务器
go run main.go

# 编译
go build -o xray-admin main.go
```

### 前端开发

```bash
# 进入前端目录
cd gin-vue-admin/web

# 安装依赖
npm install --legacy-peer-deps

# 运行开发服务器
npm run serve

# 构建生产版本
npm run build
```

### 开发环境配置

```bash
# 前端环境变量 (.env.development)
VITE_CLI_PORT = 8080
VITE_SERVER_PORT = 8888
VITE_BASE_API = /api/v1
VITE_BASE_PATH = http://127.0.0.1

# 生产环境变量 (.env.production)
VITE_CLI_PORT = 8080
VITE_SERVER_PORT = 8888
VITE_BASE_API = /api/v1
VITE_BASE_PATH = http://127.0.0.1
```

## 🔍 故障排除

### 常见问题

1. **服务启动失败**
   ```bash
   # 检查端口占用
   netstat -tlnp | grep :8888
   
   # 检查配置文件
   cat /usr/local/etc/xray-admin/config.yaml
   ```

2. **数据库连接失败**
   ```bash
   # 检查 MySQL 服务状态
   systemctl status mysql
   
   # 测试数据库连接
   mysql -u root -p -h 127.0.0.1 -P 3306
   ```

3. **流量统计不准确**
   ```bash
   # 检查 stat 服务状态
   systemctl status stat
   
   # 检查 Xray API 连接
   curl http://127.0.0.1:11111
   ```

4. **前端页面空白**
   - 清除浏览器缓存 (Ctrl+F5)
   - 检查控制台错误信息
   - 确认后端服务正常运行

### 性能优化

1. **数据库优化**
   ```sql
   -- 添加索引优化查询性能
   CREATE INDEX idx_created_at ON v2ray_stats(created_at);
   CREATE INDEX idx_server_ip ON v2ray_stats(server_ip);
   CREATE INDEX idx_tag ON v2ray_stats(tag);
   ```

2. **系统优化**
   ```bash
   # 调整文件描述符限制
   echo "* soft nofile 65535" >> /etc/security/limits.conf
   echo "* hard nofile 65535" >> /etc/security/limits.conf
   
   # 优化网络参数
   echo "net.core.rmem_max = 67108864" >> /etc/sysctl.conf
   echo "net.core.wmem_max = 67108864" >> /etc/sysctl.conf
   sysctl -p
   ```

## 📄 API 文档

### 主要 API 接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/v1/server` | GET/POST/PUT/DELETE | 服务器管理 |
| `/api/v1/binding` | GET/POST/PUT/DELETE | 绑定管理 |
| `/api/v1/stat/charts` | GET | 流量趋势数据 |
| `/api/v1/stat/rank` | GET | 流量排行数据 |
| `/api/v1/stat/list` | GET | 流量记录列表 |

### 认证方式

使用 JWT Token 认证，请求头添加：
```
Authorization: Bearer <your_jwt_token>
```

## 🤝 贡献指南

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 📝 更新日志

### v1.2.0 (2024-02-04)
- ✨ 新增流量排行榜功能
- ✨ 优化分享配置UI设计
- 🐛 修复复制配置链接功能
- 🐛 修复按钮对齐问题
- 🔧 改进缓存清理机制

### v1.1.0 (2024-01-15)
- ✨ 新增响应式设计支持
- ✨ 新增流量趋势图表
- 🐛 修复流量统计精度问题
- 🔧 优化编译脚本

### v1.0.0 (2024-01-01)
- 🎉 首次发布
- ✨ 基础服务器管理功能
- ✨ 用户绑定管理功能
- ✨ 流量统计功能

## 📜 许可证

本项目基于 [MIT License](LICENSE) 开源协议。

## 🙏 致谢

- [gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin) - 基础框架
- [Xray-core](https://github.com/XTLS/Xray-core) - 代理核心
- [Element Plus](https://element-plus.org/) - UI 组件库
- [ECharts](https://echarts.apache.org/) - 图表库

## 📞 支持

如果你在使用过程中遇到问题，可以通过以下方式获取帮助：

- 🐛 [提交 Issue](https://github.com/gongshen/xray/issues)
- 💬 [讨论区](https://github.com/gongshen/xray/discussions)
- 📧 邮件支持: support@example.com

---

⭐ 如果这个项目对你有帮助，请给个 Star 支持一下！
