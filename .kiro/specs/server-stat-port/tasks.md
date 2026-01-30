# Implementation Plan: Server Stat Port

## Overview

为代理服务器管理平台添加"统计服务端口"功能，允许每个服务器配置独立的 stat 程序端口。实现采用增量方式，先修改模型，再更新服务层，最后修改前端。

## Tasks

- [x] 1. 扩展 Server 模型
  - [x] 1.1 在 Server 结构体中添加 StatPort 字段
    - 文件: `gin-vue-admin/server/model/v2ray/server.go`
    - 添加 `StatPort int` 字段，包含 JSON、form、gorm tag
    - _Requirements: 1.1, 1.3, 1.4_
  
  - [x] 1.2 实现 GetStatPort() 方法
    - 在 Server 结构体上添加 GetStatPort() 方法
    - 当 StatPort > 0 时返回 StatPort，否则返回 global.GVA_CONFIG.STAT_PORT
    - _Requirements: 7.1, 7.2, 7.3_
  
  - [x] 1.3 编写 GetStatPort 属性测试
    - **Property 1: GetStatPort 端口回退逻辑**
    - **Validates: Requirements 1.2, 7.2, 7.3**

- [x] 2. 修改后端服务层
  - [x] 2.1 更新 ReportBinding 方法
    - 文件: `gin-vue-admin/server/service/v2ray_admin/binding.go`
    - 将 `global.GVA_CONFIG.STAT_PORT` 替换为 `srv.GetStatPort()`
    - _Requirements: 2.1, 2.2, 2.3_
  
  - [x] 2.2 更新 collectTraffic 函数
    - 文件: `gin-vue-admin/server/api/v1/job/traffic.go`
    - 将 `global.GVA_CONFIG.STAT_PORT` 替换为 `srv.GetStatPort()`
    - _Requirements: 3.1, 3.2, 3.3_
  
  - [x] 2.3 更新 collectSysInfo 函数
    - 文件: `gin-vue-admin/server/api/v1/job/traffic.go`
    - 将 `global.GVA_CONFIG.STAT_PORT` 替换为 `srv.GetStatPort()`
    - _Requirements: 4.1, 4.2, 4.3_

- [x] 3. Checkpoint - 验证后端修改
  - 确保代码编译通过，检查是否有语法错误
  - 如有问题请告知用户

- [x] 4. 修改前端表单
  - [x] 4.1 更新 serverForm.vue 添加统计端口输入框
    - 文件: `gin-vue-admin/web/src/view/v2ray_admin/server/serverForm.vue`
    - 在端口输入框后添加统计端口输入框
    - 使用 `v-model.number` 绑定 `formData.stat_port`
    - 设置 placeholder 为 "默认: 56611"
    - _Requirements: 5.1, 5.2, 5.3, 5.4_
  
  - [x] 4.2 更新 formData 初始值
    - 添加 `stat_port: 0` 到 formData 对象
    - _Requirements: 5.5_

- [x] 5. Final Checkpoint - 验证完整实现
  - 确保所有代码编译通过
  - 验证数据库迁移会自动添加 stat_port 列
  - 如有问题请告知用户

## Notes

- 每个任务引用具体的需求以便追溯
- 数据库迁移由 GORM AutoMigrate 自动处理，无需手动迁移脚本
- 属性测试验证核心正确性属性
