# Component Inventory

## Application Packages

### gin-vue-admin/server (Go 后端)
- **Purpose**: 管理平台后端服务
- **Language**: Go 1.18
- **Framework**: Gin + GORM

### gin-vue-admin/web (Vue 前端)
- **Purpose**: 管理平台前端应用
- **Language**: JavaScript/Vue 3
- **Framework**: Vue 3 + Element Plus + Vite

### stat (流量采集程序)
- **Purpose**: 代理服务器流量采集
- **Language**: Go 1.17
- **Framework**: fasthttp

## Backend Module Inventory

### 核心业务模块 (保留)

| 模块 | 路径 | 用途 |
|------|------|------|
| v2ray_admin | api/v1/v2ray_admin/ | 管理员代理管理 API |
| v2ray | api/v1/v2ray/ | 用户代理 API |
| v2ray_admin | service/v2ray_admin/ | 代理管理业务逻辑 |
| v2ray | model/v2ray/ | 代理数据模型 |
| v2ray_admin | router/v2ray_admin/ | 管理员路由 |
| v2ray | router/v2ray/ | 用户路由 |

### 系统模块 (保留)

| 模块 | 路径 | 用途 |
|------|------|------|
| system | api/v1/system/ | 系统管理 API |
| system | service/system/ | 系统业务逻辑 |
| system | model/system/ | 系统数据模型 |
| system | router/system/ | 系统路由 |
| example | api/v1/example/ | 示例 API |
| example | service/example/ | 示例业务逻辑 |

### 待删除模块 (自动代码生成)

| 模块 | 路径 | 用途 | 状态 |
|------|------|------|------|
| AutoCode | api/v1/system/sys_auto_code.go | 代码生成 API | 待删除 |
| AutoCodeHistory | api/v1/system/sys_auto_code_history.go | 代码历史 API | 待删除 |
| AutoCode | service/system/sys_auto_code.go | 代码生成服务 | 待删除 |
| AutoCodeHistory | service/system/sys_autocode_history.go | 代码历史服务 | 待删除 |
| AutoCode | router/system/sys_auto_code.go | 代码生成路由 | 待删除 |
| AutoCodeHistory | router/system/sys_auto_code_history.go | 代码历史路由 | 待删除 |
| AutoCode | model/system/sys_auto_code.go | 代码生成模型 | 待删除 |
| AutoCodeHistory | model/system/sys_autocode_history.go | 代码历史模型 | 待删除 |
| AST Utils | utils/ast/ | AST 代码注入工具 | 待删除 |
| Templates | resource/autocode_template/ | 代码模板 | 待删除 |
| Plug Templates | resource/plug_template/ | 插件模板 | 待删除 |

## Frontend Module Inventory

### 核心业务页面 (保留)

| 模块 | 路径 | 用途 |
|------|------|------|
| v2ray_admin | view/v2ray_admin/ | 管理员代理管理页面 |
| v2ray | view/v2ray/ | 用户代理页面 |
| superAdmin | view/superAdmin/ | 超级管理员页面 |
| dashboard | view/dashboard/ | 仪表盘 |
| login | view/login/ | 登录页面 |
| layout | view/layout/ | 布局组件 |

### 待删除页面 (系统工具)

| 模块 | 路径 | 用途 | 状态 |
|------|------|------|------|
| autoCode | view/systemTools/autoCode/ | 代码生成器 | 待删除 |
| autoCodeAdmin | view/systemTools/autoCodeAdmin/ | 代码管理 | 待删除 |
| autoPkg | view/systemTools/autoPkg/ | 自动化 Package | 待删除 |
| autoPlug | view/systemTools/autoPlug/ | 自动插件 | 待删除 |
| formCreate | view/systemTools/formCreate/ | 表单生成器 | 待删除 |
| installPlugin | view/systemTools/installPlugin/ | 插件安装 | 待删除 |

### 待删除 API 文件

| 文件 | 用途 | 状态 |
|------|------|------|
| api/autoCode.js | 自动代码 API | 待删除 |

## Shared Packages

| 包 | 路径 | 用途 |
|------|------|------|
| global | server/global/ | 全局变量 |
| config | server/config/ | 配置结构 |
| middleware | server/middleware/ | 中间件 |
| utils | server/utils/ | 工具函数 |
| core | server/core/ | 核心启动 |

## Total Count

| 类型 | 数量 |
|------|------|
| **总包数** | 3 |
| Application | 3 (server, web, stat) |
| Infrastructure | 0 |
| Shared | 5 (global, config, middleware, utils, core) |
| Test | 0 (无独立测试包) |

## 待删除文件统计

| 类型 | 数量 |
|------|------|
| 后端 Go 文件 | ~15 |
| 前端 Vue/JS 文件 | ~10 |
| 模板目录 | 2 |
| **总计** | ~27 文件/目录 |
