# Dependencies

## Internal Dependencies

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          组件依赖关系                                        │
│                                                                             │
│   ┌─────────────────┐                                                       │
│   │  gin-vue-admin  │                                                       │
│   │     /web        │                                                       │
│   └────────┬────────┘                                                       │
│            │ HTTP API                                                       │
│            ▼                                                                │
│   ┌─────────────────┐         ┌─────────────────┐                          │
│   │  gin-vue-admin  │ ──────► │     MySQL       │                          │
│   │    /server      │  GORM   │                 │                          │
│   └────────┬────────┘         └─────────────────┘                          │
│            │                                                                │
│            │ HTTP                                                           │
│            ▼                                                                │
│   ┌─────────────────┐         ┌─────────────────┐                          │
│   │      stat       │ ──────► │    xray-core    │                          │
│   │   (各服务器)     │  gRPC   │   (代理服务)     │                          │
│   └─────────────────┘         └─────────────────┘                          │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### gin-vue-admin/web → gin-vue-admin/server
- **Type**: Runtime
- **Reason**: 前端通过 HTTP API 调用后端服务

### gin-vue-admin/server → MySQL
- **Type**: Runtime
- **Reason**: 数据持久化

### gin-vue-admin/server → Redis
- **Type**: Runtime
- **Reason**: JWT 黑名单、缓存

### gin-vue-admin/server → stat
- **Type**: Runtime
- **Reason**: 采集流量数据、下发配置

### stat → xray-core
- **Type**: Runtime
- **Reason**: 通过 gRPC 获取流量统计

## External Dependencies

### Backend (gin-vue-admin/server)

| 依赖 | 版本 | 用途 | License |
|------|------|------|---------|
| gin-gonic/gin | 1.7.0 | Web 框架 | MIT |
| gorm.io/gorm | 1.23.4 | ORM | MIT |
| casbin/casbin | 2.51.0 | RBAC | Apache-2.0 |
| golang-jwt/jwt | 4.3.0 | JWT | MIT |
| spf13/viper | 1.7.0 | 配置管理 | MIT |
| uber-go/zap | 1.16.0 | 日志 | MIT |
| go-redis/redis | 8.11.4 | Redis 客户端 | BSD-2 |
| xtls/xray-core | 1.8.0 | 代理核心 | MPL-2.0 |
| sashabaranov/go-openai | 1.5.7 | OpenAI API | Apache-2.0 |
| shirou/gopsutil | 3.22.5 | 系统信息 | BSD-3 |
| otiai10/copy | 1.7.0 | 文件复制 | MIT |

### Backend (stat)

| 依赖 | 版本 | 用途 | License |
|------|------|------|---------|
| valyala/fasthttp | 1.45.0 | HTTP 服务器 | MIT |
| sirupsen/logrus | 1.9.0 | 日志 | MIT |
| xtls/xray-core | 1.4.2 | gRPC 客户端 | MPL-2.0 |
| grpc/grpc-go | 1.49.0 | gRPC | Apache-2.0 |

### Frontend (gin-vue-admin/web)

| 依赖 | 版本 | 用途 | License |
|------|------|------|---------|
| vue | 3.2.25 | 前端框架 | MIT |
| element-plus | 2.2.30 | UI 组件库 | MIT |
| pinia | 2.0.9 | 状态管理 | MIT |
| vue-router | 4.0.0 | 路由 | MIT |
| axios | 1.8.3 | HTTP 客户端 | MIT |
| echarts | 5.4.2 | 图表 | Apache-2.0 |
| lodash | 4.17.21 | 工具库 | MIT |
| marked | 4.3.0 | Markdown 解析 | MIT |
| qrcode | 1.5.1 | 二维码生成 | MIT |

### Dev Dependencies (Frontend)

| 依赖 | 版本 | 用途 | License |
|------|------|------|---------|
| vite | 4.4.9 | 构建工具 | MIT |
| @vitejs/plugin-vue | 5.2.1 | Vue 插件 | MIT |
| sass | 1.54.0 | CSS 预处理 | MIT |
| eslint | 6.7.2 | 代码检查 | MIT |

## 待删除依赖

以下依赖主要用于自动代码生成功能，删除相关功能后可能不再需要：

| 依赖 | 用途 | 状态 |
|------|------|------|
| otiai10/copy | 文件复制 (代码生成) | 评估是否保留 |
| utils/ast | AST 代码注入 | 待删除 |

## 依赖更新建议

| 依赖 | 当前版本 | 建议 |
|------|----------|------|
| gin | 1.7.0 | 可升级到 1.9.x |
| gorm | 1.23.4 | 可升级到 1.25.x |
| vue | 3.2.25 | 可升级到 3.4.x |
| vite | 4.4.9 | 可升级到 5.x |
| element-plus | 2.2.30 | 可升级到 2.4.x |
