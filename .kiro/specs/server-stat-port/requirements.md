# Requirements Document

## Introduction

本功能为代理服务器管理平台添加"统计服务端口"（Stat Port）配置能力。当前系统通过全局配置 `config.yaml` 中的 `stat_port: 56611` 访问所有代理服务器的 stat 程序，无法满足不同服务器使用不同端口的需求。本功能允许每个服务器配置独立的统计服务端口，使后端能够通过该端口与不同服务器的 stat 程序通信。

## Glossary

- **Server**: 代理服务器实体，存储于 `v2ray_server` 表
- **Stat_Port**: 统计服务端口，用于与服务器上的 stat 程序通信
- **Stat_Program**: 运行在代理服务器上的统计程序，提供流量统计和系统信息接口
- **Binding_Service**: 绑定服务，负责上报用户绑定配置到代理服务器
- **Traffic_Collector**: 流量收集器，定时从各服务器收集流量统计数据
- **SysInfo_Collector**: 系统信息收集器，定时从各服务器收集系统资源使用情况
- **Global_Config**: 全局配置，存储于 `config.yaml` 的默认配置值

## Requirements

### Requirement 1: Server 模型扩展

**User Story:** As a 系统管理员, I want to 为每个服务器配置独立的统计服务端口, so that 不同服务器可以使用不同的 stat 程序端口。

#### Acceptance Criteria

1. THE Server 模型 SHALL 包含 `StatPort` 字段，类型为 `int`，用于存储统计服务端口
2. WHEN StatPort 字段值为 0 或未设置时, THE System SHALL 使用全局配置 `global.GVA_CONFIG.STAT_PORT` 作为默认值
3. THE StatPort 字段 SHALL 支持 JSON 序列化，字段名为 `stat_port`
4. THE StatPort 字段 SHALL 映射到数据库列 `stat_port`

### Requirement 2: 绑定上报功能适配

**User Story:** As a 系统管理员, I want to 绑定上报功能使用服务器独立的统计端口, so that 配置能够正确推送到各服务器的 stat 程序。

#### Acceptance Criteria

1. WHEN ReportBinding 方法执行时, THE Binding_Service SHALL 优先使用 Server.StatPort 作为通信端口
2. IF Server.StatPort 为 0, THEN THE Binding_Service SHALL 回退使用 global.GVA_CONFIG.STAT_PORT
3. THE ReportBinding 方法 SHALL 构建正确的 HTTP 请求 URI: `http://{server.Ip}:{effectivePort}/conf/update`

### Requirement 3: 流量收集功能适配

**User Story:** As a 系统管理员, I want to 流量收集功能使用服务器独立的统计端口, so that 能够从各服务器正确获取流量统计数据。

#### Acceptance Criteria

1. WHEN collectTraffic 方法执行时, THE Traffic_Collector SHALL 优先使用 Server.StatPort 作为通信端口
2. IF Server.StatPort 为 0, THEN THE Traffic_Collector SHALL 回退使用 global.GVA_CONFIG.STAT_PORT
3. THE collectTraffic 方法 SHALL 构建正确的 HTTP 请求 URI: `http://{server.Ip}:{effectivePort}/stat/traffic`

### Requirement 4: 系统信息收集功能适配

**User Story:** As a 系统管理员, I want to 系统信息收集功能使用服务器独立的统计端口, so that 能够从各服务器正确获取系统资源使用情况。

#### Acceptance Criteria

1. WHEN collectSysInfo 方法执行时, THE SysInfo_Collector SHALL 优先使用 Server.StatPort 作为通信端口
2. IF Server.StatPort 为 0, THEN THE SysInfo_Collector SHALL 回退使用 global.GVA_CONFIG.STAT_PORT
3. THE collectSysInfo 方法 SHALL 构建正确的 HTTP 请求 URI: `http://{server.Ip}:{effectivePort}/stat/sysinfo`

### Requirement 5: 前端表单扩展

**User Story:** As a 系统管理员, I want to 在服务器新增/编辑页面配置统计服务端口, so that 可以通过界面管理每个服务器的 stat 端口。

#### Acceptance Criteria

1. THE serverForm.vue SHALL 包含"统计服务端口"输入框
2. THE 统计服务端口输入框 SHALL 显示标签"统计端口:"
3. THE 统计服务端口输入框 SHALL 接受数字类型输入
4. WHEN 统计服务端口为空或 0 时, THE 前端 SHALL 显示占位提示"默认: 56611"
5. THE formData SHALL 包含 `stat_port` 字段，默认值为 0

### Requirement 6: 数据库迁移

**User Story:** As a 系统管理员, I want to 数据库自动添加 stat_port 字段, so that 现有系统可以平滑升级。

#### Acceptance Criteria

1. THE System SHALL 通过 GORM AutoMigrate 自动添加 `stat_port` 列到 `v2ray_server` 表
2. THE stat_port 列 SHALL 允许 NULL 值或默认为 0
3. WHEN 现有服务器记录的 stat_port 为 NULL 或 0 时, THE System SHALL 使用全局默认端口

### Requirement 7: 端口获取辅助函数

**User Story:** As a 开发者, I want to 有统一的辅助函数获取服务器的有效统计端口, so that 代码逻辑清晰且易于维护。

#### Acceptance Criteria

1. THE Server 模型 SHALL 提供 `GetStatPort()` 方法返回有效的统计端口
2. WHEN Server.StatPort 大于 0 时, THE GetStatPort() 方法 SHALL 返回 Server.StatPort
3. WHEN Server.StatPort 为 0 时, THE GetStatPort() 方法 SHALL 返回 global.GVA_CONFIG.STAT_PORT
