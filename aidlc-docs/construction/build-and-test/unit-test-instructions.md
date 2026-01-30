# Unit Test Instructions - 代理服务器管理平台

## 概述
本文档提供项目单元测试执行指南。

---

## 1. 后端单元测试 (gin-vue-admin/server)

### 现有测试文件
```
gin-vue-admin/server/utils/
├── human_duration_test.go    # 时间格式化测试
├── validator_test.go         # 验证器测试
└── timer/
    └── timed_task_test.go    # 定时任务测试
```

### 执行测试

```bash
# 进入后端目录
cd gin-vue-admin/server

# 运行所有测试
go test ./...

# 运行特定包测试
go test ./utils/...

# 带详细输出
go test -v ./...

# 带覆盖率
go test -cover ./...
```

### 预期结果
```
ok      gin-vue-admin/server/utils      0.XXXs
ok      gin-vue-admin/server/utils/timer 0.XXXs
```

---

## 2. 前端单元测试 (gin-vue-admin/web)

### 测试框架
项目当前未配置前端单元测试框架。

### 建议配置 (可选)
如需添加前端测试，建议使用 Vitest：

```bash
# 安装 Vitest
npm install -D vitest @vue/test-utils happy-dom

# 在 package.json 添加脚本
# "test": "vitest run"
# "test:watch": "vitest"
```

---

## 3. Stat 程序单元测试 (stat)

### 现有测试
stat 程序当前无单元测试文件。

### 执行测试 (如有)
```bash
cd stat
go test ./...
```

---

## 测试执行清单

| 组件 | 命令 | 状态 |
|------|------|------|
| 后端 | `go test ./...` | ✅ 有测试 |
| 前端 | N/A | ⚠️ 无测试框架 |
| Stat | `go test ./...` | ⚠️ 无测试文件 |

---

## 测试覆盖率报告

### 生成覆盖率报告
```bash
cd gin-vue-admin/server

# 生成覆盖率文件
go test -coverprofile=coverage.out ./...

# 查看覆盖率摘要
go tool cover -func=coverage.out

# 生成 HTML 报告
go tool cover -html=coverage.out -o coverage.html
```

---

## 注意事项

1. **数据库依赖**: 部分测试可能需要数据库连接，确保配置正确
2. **环境变量**: 检查测试是否需要特定环境变量
3. **并发测试**: 使用 `-race` 标志检测竞态条件
   ```bash
   go test -race ./...
   ```
