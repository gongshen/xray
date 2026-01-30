# Code Structure

## Build System

- **Type**: Go Modules + npm/Vite
- **Configuration**: 
  - Backend: `gin-vue-admin/server/go.mod`
  - Frontend: `gin-vue-admin/web/package.json`
  - Stat: `stat/go.mod`

## Backend Structure (gin-vue-admin/server)

```
server/
├── api/v1/                    # API 控制器层
│   ├── enter.go               # API 组入口
│   ├── system/                # 系统管理 API
│   ├── example/               # 示例 API
│   ├── v2ray/                 # 用户代理 API
│   └── v2ray_admin/           # 管理员代理 API
├── config/                    # 配置结构定义
├── core/                      # 核心启动逻辑
├── global/                    # 全局变量
├── initialize/                # 初始化逻辑
├── middleware/                # 中间件 (JWT, Casbin, CORS)
├── model/                     # 数据模型
│   ├── common/                # 通用请求/响应
│   ├── system/                # 系统模型
│   └── v2ray/                 # 代理相关模型
├── router/                    # 路由定义
│   ├── system/                # 系统路由
│   ├── v2ray/                 # 用户代理路由
│   └── v2ray_admin/           # 管理员代理路由
├── service/                   # 业务逻辑层
│   ├── system/                # 系统服务
│   └── v2ray_admin/           # 代理管理服务
├── source/                    # 初始化数据源
├── utils/                     # 工具函数
└── resource/                  # 静态资源和模板
```

## Frontend Structure (gin-vue-admin/web)

```
web/src/
├── api/                       # API 调用
│   ├── server.js              # 服务器管理 API
│   ├── stat.js                # 流量统计 API
│   ├── binding.js             # 绑定管理 API
│   ├── autoCode.js            # [待删除] 自动代码 API
│   └── ...
├── view/                      # 页面组件
│   ├── v2ray/                 # 用户代理页面
│   ├── v2ray_admin/           # 管理员代理页面
│   ├── superAdmin/            # 超级管理员页面
│   ├── systemTools/           # [待删除] 系统工具页面
│   │   ├── autoCode/          # [待删除] 代码生成器
│   │   ├── autoCodeAdmin/     # [待删除] 代码管理
│   │   ├── autoPkg/           # [待删除] 自动化 Package
│   │   ├── autoPlug/          # [待删除] 自动插件
│   │   ├── formCreate/        # [待删除] 表单生成器
│   │   └── installPlugin/     # [待删除] 插件安装
│   └── ...
├── components/                # 公共组件
├── pinia/                     # 状态管理
├── router/                    # 路由配置
├── style/                     # 样式文件
└── utils/                     # 工具函数
```

## Stat Program Structure

```
stat/
├── main.go                    # 入口
├── business/                  # 业务逻辑
│   ├── stat.go                # 流量采集
│   ├── config.go              # 配置更新
│   └── command.go             # 系统命令
├── conn/                      # xray gRPC 连接
├── server/                    # HTTP 服务器
└── utils/                     # 工具函数
```

## Key Files Inventory

### 核心业务文件 (保留)
- `server/model/v2ray/server.go` - 服务器模型
- `server/model/v2ray/stat.go` - 流量统计模型
- `server/model/v2ray/binding.go` - 用户绑定模型
- `server/service/v2ray_admin/server.go` - 服务器管理服务
- `server/service/v2ray_admin/stat.go` - 流量统计服务
- `server/service/v2ray_admin/binding.go` - 绑定管理服务
- `server/api/v1/v2ray_admin/*.go` - 管理员 API
- `server/api/v1/v2ray/*.go` - 用户 API
- `server/router/v2ray_admin/*.go` - 管理员路由
- `server/router/v2ray/*.go` - 用户路由

### 待删除文件 (自动代码生成相关)
- `server/router/system/sys_auto_code.go` - 自动代码路由
- `server/router/system/sys_auto_code_history.go` - 代码历史路由
- `server/service/system/sys_auto_code.go` - 自动代码服务
- `server/service/system/sys_auto_code_history.go` - 代码历史服务
- `server/service/system/sys_auto_code_*.go` - 各数据库适配
- `server/api/v1/system/sys_auto_code.go` - 自动代码 API
- `server/api/v1/system/sys_auto_code_history.go` - 代码历史 API
- `server/model/system/sys_auto_code.go` - 自动代码模型
- `server/model/system/sys_autocode_history.go` - 代码历史模型
- `server/resource/autocode_template/` - 代码模板目录
- `server/resource/plug_template/` - 插件模板目录
- `server/utils/ast/` - AST 工具 (代码生成用)
- `web/src/api/autoCode.js` - 前端自动代码 API
- `web/src/view/systemTools/autoCode/` - 代码生成器页面
- `web/src/view/systemTools/autoCodeAdmin/` - 代码管理页面
- `web/src/view/systemTools/autoPkg/` - 自动 Package 页面
- `web/src/view/systemTools/autoPlug/` - 自动插件页面
- `web/src/view/systemTools/formCreate/` - 表单生成器页面
- `web/src/view/systemTools/installPlugin/` - 插件安装页面

## Design Patterns

### Repository Pattern
- **Location**: `service/` 层
- **Purpose**: 封装数据访问逻辑
- **Implementation**: 每个业务模块有对应的 Service 结构体

### Middleware Pattern
- **Location**: `middleware/`
- **Purpose**: 请求处理链
- **Implementation**: JWT 认证、Casbin 权限、日志记录

### Dependency Injection
- **Location**: `global/global.go`
- **Purpose**: 全局依赖管理
- **Implementation**: 全局变量存储 DB、Config、Logger 等

## Critical Dependencies

### Backend
| 依赖 | 版本 | 用途 |
|------|------|------|
| gin | 1.7.0 | Web 框架 |
| gorm | 1.23.4 | ORM |
| casbin | 2.51.0 | RBAC 权限 |
| jwt-go | 4.3.0 | JWT 认证 |
| xray-core | 1.8.0 | 代理核心 |
| zap | 1.16.0 | 日志 |
| viper | 1.7.0 | 配置管理 |

### Frontend
| 依赖 | 版本 | 用途 |
|------|------|------|
| vue | 3.2.25 | 前端框架 |
| element-plus | 2.2.30 | UI 组件库 |
| pinia | 2.0.9 | 状态管理 |
| axios | 1.8.3 | HTTP 客户端 |
| echarts | 5.4.2 | 图表 |
| vite | 4.4.9 | 构建工具 |
