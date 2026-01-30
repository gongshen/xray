# Technology Stack

## Programming Languages

| 语言 | 版本 | 用途 |
|------|------|------|
| Go | 1.18 | 后端服务 (gin-vue-admin/server) |
| Go | 1.17 | 流量采集程序 (stat) |
| JavaScript | ES6+ | 前端应用 |
| Vue | 3.2.25 | 前端框架 |
| SCSS | - | 样式 |

## Frameworks

### Backend Frameworks

| 框架 | 版本 | 用途 |
|------|------|------|
| Gin | 1.7.0 | Web 框架 |
| GORM | 1.23.4 | ORM |
| Casbin | 2.51.0 | RBAC 权限控制 |
| Viper | 1.7.0 | 配置管理 |
| Zap | 1.16.0 | 日志 |
| fasthttp | 1.45.0 | HTTP 服务器 (stat) |

### Frontend Frameworks

| 框架 | 版本 | 用途 |
|------|------|------|
| Vue | 3.2.25 | 前端框架 |
| Element Plus | 2.2.30 | UI 组件库 |
| Pinia | 2.0.9 | 状态管理 |
| Vue Router | 4.0.0 | 路由 |
| Axios | 1.8.3 | HTTP 客户端 |
| ECharts | 5.4.2 | 图表 |

## Infrastructure

| 服务 | 用途 |
|------|------|
| MySQL | 主数据库 |
| Redis | JWT 黑名单、缓存 |
| xray-core | 代理服务 |
| Docker | 容器化部署 |
| Kubernetes | 容器编排 (可选) |

## Build Tools

| 工具 | 版本 | 用途 |
|------|------|------|
| Go Modules | - | Go 依赖管理 |
| npm | - | Node.js 包管理 |
| Vite | 4.4.9 | 前端构建工具 |
| Docker | - | 容器构建 |

## Testing Tools

| 工具 | 版本 | 用途 |
|------|------|------|
| Go testing | - | Go 单元测试 |
| testify | 1.8.2 | Go 测试断言 |

## Security

| 技术 | 用途 |
|------|------|
| JWT | 用户认证 |
| Casbin | RBAC 权限控制 |
| bcrypt | 密码加密 |
| CORS | 跨域控制 |

## External Services (可选)

| 服务 | 用途 |
|------|------|
| 七牛云 OSS | 文件存储 |
| 阿里云 OSS | 文件存储 |
| 腾讯云 COS | 文件存储 |
| AWS S3 | 文件存储 |
| 华为云 OBS | 文件存储 |

## Development Tools

| 工具 | 用途 |
|------|------|
| Swagger | API 文档 |
| ESLint | JavaScript 代码检查 |
| Goland/VSCode | IDE |
