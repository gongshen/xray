# Requirements Document

## Intent Analysis Summary

### User Request
用户希望对代理服务器管理平台进行代码清理和优化：
1. 优化不合理的前端或后端代码
2. 删掉没有用到的废代码
3. 删除系统工具中的自动化Package、自动化代码管理、代码生成器、表单生成器、万用表格功能

### Request Type
- **Refactoring** + **Code Cleanup**

### Scope Estimate
- **Multiple Components**: 涉及前端和后端多个模块

### Complexity Estimate
- **Moderate**: 需要识别和删除多个相关联的模块，同时确保不破坏核心功能

---

## Functional Requirements

### FR-1: 删除自动代码生成功能
**优先级**: 高

删除以下功能模块：
- FR-1.1: 删除代码生成器 (autoCode) - 用于自动生成 CRUD 代码
- FR-1.2: 删除代码管理 (autoCodeAdmin) - 用于管理生成的代码历史
- FR-1.3: 删除自动化 Package (autoPkg) - 用于创建新的业务包
- FR-1.4: 删除自动插件 (autoPlug) - 用于创建插件模板
- FR-1.5: 删除插件安装 (installPlugin) - 用于安装第三方插件

### FR-2: 删除表单生成器功能
**优先级**: 高

- FR-2.1: 删除表单生成器页面 (formCreate)
- FR-2.2: 删除相关的前端组件和 API

### FR-3: 清理废弃代码
**优先级**: 高

- FR-3.1: 删除 `gin-vue-admin/rm_file/` 目录 (历史删除文件备份)
- FR-3.2: 清理未使用的导入和引用

### FR-4: 更新路由和入口文件
**优先级**: 高

- FR-4.1: 从 `router.go` 中移除 AutoCode 相关路由注册
- FR-4.2: 从 `enter.go` 文件中移除相关模块引用
- FR-4.3: 更新前端路由配置，移除相关菜单

### FR-5: 保持核心功能完整
**优先级**: 关键

确保以下核心功能不受影响：
- FR-5.1: 用户认证和授权 (JWT + Casbin)
- FR-5.2: 代理服务器管理 (v2ray_admin)
- FR-5.3: 用户绑定管理
- FR-5.4: 流量统计功能
- FR-5.5: 系统管理功能 (用户、角色、菜单、API 管理)

---

## Non-Functional Requirements

### NFR-1: 代码质量
- NFR-1.1: 删除后代码应能正常编译
- NFR-1.2: 删除后前端应能正常构建
- NFR-1.3: 不引入新的代码警告或错误

### NFR-2: 向后兼容
- NFR-2.1: 现有数据库数据不受影响
- NFR-2.2: 现有 API 接口保持兼容
- NFR-2.3: 用户登录状态不受影响

### NFR-3: 可维护性
- NFR-3.1: 删除后代码结构更清晰
- NFR-3.2: 减少不必要的依赖

---

## Scope Definition

### In Scope (范围内)

| 组件 | 操作 |
|------|------|
| 后端 AutoCode 相关 | 删除 |
| 后端 AutoCodeHistory 相关 | 删除 |
| 后端 AST 工具 | 删除 |
| 后端代码模板 | 删除 |
| 前端 autoCode 页面 | 删除 |
| 前端 autoCodeAdmin 页面 | 删除 |
| 前端 autoPkg 页面 | 删除 |
| 前端 autoPlug 页面 | 删除 |
| 前端 formCreate 页面 | 删除 |
| 前端 installPlugin 页面 | 删除 |
| 前端 autoCode API | 删除 |
| rm_file 目录 | 删除 |
| 入口文件引用 | 更新 |
| 路由注册 | 更新 |

### Out of Scope (范围外)

| 组件 | 说明 |
|------|------|
| 代理管理功能 | 保持不变 |
| 用户管理功能 | 保持不变 |
| 流量统计功能 | 保持不变 |
| stat 程序 | 保持不变 |
| 数据库结构 | 保持不变 |
| 依赖升级 | 本次不处理 |
| 添加测试 | 本次不处理 |

---

## Acceptance Criteria

### AC-1: 功能删除验证
- [ ] 后端编译成功，无 AutoCode 相关代码
- [ ] 前端构建成功，无 systemTools 下的删除页面
- [ ] 访问已删除功能的 URL 返回 404

### AC-2: 核心功能验证
- [ ] 用户登录正常
- [ ] 代理服务器管理正常
- [ ] 流量统计正常
- [ ] 用户绑定管理正常

### AC-3: 代码清洁度
- [ ] 无未使用的导入
- [ ] 无孤立的文件引用
- [ ] rm_file 目录已删除

---

## File Deletion List

### Backend Files to Delete

```
gin-vue-admin/server/
├── api/v1/system/
│   ├── sys_auto_code.go
│   └── sys_auto_code_history.go
├── service/system/
│   ├── sys_auto_code.go
│   ├── sys_autocode_history.go
│   ├── sys_auto_code_interface.go
│   ├── sys_auto_code_mysql.go
│   ├── sys_auto_code_pgsql.go
│   ├── sys_auto_code_mssql.go
│   └── sys_auto_code_oracle.go
├── router/system/
│   ├── sys_auto_code.go
│   └── sys_auto_code_history.go
├── model/system/
│   ├── sys_auto_code.go
│   └── sys_autocode_history.go
├── resource/
│   ├── autocode_template/  (整个目录)
│   └── plug_template/      (整个目录)
└── utils/ast/              (整个目录)
```

### Frontend Files to Delete

```
gin-vue-admin/web/src/
├── api/
│   └── autoCode.js
└── view/systemTools/
    ├── autoCode/           (整个目录)
    ├── autoCodeAdmin/      (整个目录)
    ├── autoPkg/            (整个目录)
    ├── autoPlug/           (整个目录)
    ├── formCreate/         (整个目录)
    └── installPlugin/      (整个目录)
```

### Other Files to Delete

```
gin-vue-admin/rm_file/      (整个目录)
```

---

## Files to Modify

### Backend Files

| 文件 | 修改内容 |
|------|----------|
| `api/v1/enter.go` | 移除 AutoCode 相关导入和引用 |
| `service/enter.go` | 移除 AutoCode 相关导入和引用 |
| `router/system/enter.go` | 移除 AutoCode 路由结构体 |
| `initialize/router.go` | 移除 AutoCode 路由注册调用 |
| `config/auto_code.go` | 评估是否删除或保留 |

### Frontend Files

| 文件 | 修改内容 |
|------|----------|
| 路由配置 | 移除相关路由 |
| 菜单配置 | 移除相关菜单项 |

---

## Risk Assessment

| 风险 | 可能性 | 影响 | 缓解措施 |
|------|--------|------|----------|
| 删除后编译失败 | 中 | 高 | 逐步删除，每步验证 |
| 遗漏引用导致运行时错误 | 中 | 中 | 全局搜索相关引用 |
| 误删核心功能代码 | 低 | 高 | 仔细审查删除列表 |
| 数据库菜单数据残留 | 低 | 低 | 提供清理 SQL |
