# Integration Test Instructions - 代理服务器管理平台

## 概述
本文档提供系统集成测试指南，验证各组件之间的交互。

---

## 1. 后端 API 集成测试

### 测试环境准备

1. **启动数据库**
   ```bash
   # MySQL (根据 config.yaml 配置)
   # 确保数据库服务运行中
   ```

2. **启动 Redis** (如配置)
   ```bash
   # 确保 Redis 服务运行中
   ```

3. **配置文件检查**
   ```bash
   # 检查 gin-vue-admin/server/config.yaml
   # 确保数据库连接信息正确
   ```

### 启动后端服务
```bash
cd gin-vue-admin/server
go run main.go
# 或使用编译后的可执行文件
./server
```

### API 测试用例

#### 1.1 健康检查
```bash
curl http://localhost:8888/health
# 预期: 200 OK
```

#### 1.2 登录接口
```bash
curl -X POST http://localhost:8888/base/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456","captcha":"","captchaId":""}'
# 预期: 返回 token
```

#### 1.3 获取用户信息 (需要 token)
```bash
curl http://localhost:8888/user/getUserInfo \
  -H "x-token: <your-token>"
# 预期: 返回用户信息
```

---

## 2. 前后端集成测试

### 测试步骤

1. **启动后端服务**
   ```bash
   cd gin-vue-admin/server
   go run main.go
   ```

2. **启动前端开发服务器**
   ```bash
   cd gin-vue-admin/web
   npm run serve
   ```

3. **访问前端页面**
   - 打开浏览器访问 `http://localhost:8080`
   - 使用默认账号登录: admin / 123456

### 功能验证清单

| 功能 | 测试步骤 | 预期结果 |
|------|----------|----------|
| 登录 | 输入账号密码点击登录 | 成功进入首页 |
| 用户管理 | 访问用户管理页面 | 显示用户列表 |
| 角色管理 | 访问角色管理页面 | 显示角色列表 |
| 菜单管理 | 访问菜单管理页面 | 显示菜单树 |
| V2Ray 管理 | 访问代理服务器页面 | 显示服务器列表 |

---

## 3. Stat 与后端集成测试

### 测试环境

1. **配置 Stat 程序**
   ```bash
   # 编辑 stat 配置，指向后端 API
   ```

2. **启动后端服务**
   ```bash
   cd gin-vue-admin/server
   go run main.go
   ```

3. **启动 Stat 程序**
   ```bash
   cd stat
   go run main.go
   ```

### 验证点

| 测试项 | 验证方法 | 预期结果 |
|--------|----------|----------|
| 流量上报 | 检查后端日志 | 收到流量数据 |
| 服务器状态 | 检查数据库 | 状态已更新 |
| 心跳检测 | 监控连接 | 保持连接 |

---

## 4. 删除功能验证

### 验证已删除功能不可访问

| 已删除功能 | 测试方法 | 预期结果 |
|------------|----------|----------|
| 自动化代码 | 访问 /systemTools/autoCode | 404 或无菜单 |
| 代码生成器 | 访问 /systemTools/autoCodeAdmin | 404 或无菜单 |
| 表单生成器 | 访问 /systemTools/formCreate | 404 或无菜单 |
| 万用表格 | 检查菜单 | 无此选项 |

### API 验证
```bash
# 验证已删除的 API 不可访问
curl http://localhost:8888/autoCode/getDB
# 预期: 404 Not Found
```

---

## 测试报告模板

```markdown
## 集成测试报告

**测试日期**: YYYY-MM-DD
**测试人员**: 
**环境**: 

### 测试结果

| 测试项 | 状态 | 备注 |
|--------|------|------|
| 后端启动 | ✅/❌ | |
| 前端启动 | ✅/❌ | |
| 登录功能 | ✅/❌ | |
| 用户管理 | ✅/❌ | |
| 已删除功能 | ✅/❌ | |

### 问题记录
- 
```
