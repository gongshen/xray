# Code Generation Plan - 代码清理

## Unit Context
- **Unit Name**: code-cleanup
- **Purpose**: 删除自动代码生成相关功能，清理废弃代码
- **Risk Level**: Low
- **Dependencies**: None

---

## Code Generation Steps

### Phase 1: 后端文件删除

- [x] Step 1.1: 删除后端 AutoCode API 文件
  - `gin-vue-admin/server/api/v1/system/sys_auto_code.go` ✓
  - `gin-vue-admin/server/api/v1/system/sys_auto_code_history.go` ✓

- [x] Step 1.2: 删除后端 AutoCode Service 文件
  - `gin-vue-admin/server/service/system/sys_auto_code.go` ✓
  - `gin-vue-admin/server/service/system/sys_autocode_history.go` ✓
  - `gin-vue-admin/server/service/system/sys_auto_code_interface.go` ✓
  - `gin-vue-admin/server/service/system/sys_auto_code_mysql.go` ✓
  - `gin-vue-admin/server/service/system/sys_auto_code_pgsql.go` ✓
  - `gin-vue-admin/server/service/system/sys_auto_code_mssql.go` ✓
  - `gin-vue-admin/server/service/system/sys_auto_code_oracle.go` ✓

- [x] Step 1.3: 删除后端 AutoCode Router 文件
  - `gin-vue-admin/server/router/system/sys_auto_code.go` ✓
  - `gin-vue-admin/server/router/system/sys_auto_code_history.go` ✓

- [x] Step 1.4: 删除后端 AutoCode Model 文件
  - `gin-vue-admin/server/model/system/sys_auto_code.go` ✓
  - `gin-vue-admin/server/model/system/sys_autocode_history.go` ✓

- [x] Step 1.5: 删除后端模板目录
  - `gin-vue-admin/server/resource/autocode_template/` (整个目录) ✓
  - `gin-vue-admin/server/resource/plug_template/` (整个目录) ✓

- [x] Step 1.6: 删除后端 AST 工具目录
  - `gin-vue-admin/server/utils/ast/` (整个目录) ✓

### Phase 2: 后端入口文件修改

- [x] Step 2.1: 修改 `gin-vue-admin/server/api/v1/system/enter.go`
  - 移除 AutoCodeApi 和 AutoCodeHistoryApi 引用 ✓

- [x] Step 2.2: 修改 `gin-vue-admin/server/service/system/enter.go`
  - 移除 AutoCodeService 和 AutoCodeHistoryService 引用 ✓

- [x] Step 2.3: 修改 `gin-vue-admin/server/router/system/enter.go`
  - 移除 AutoCodeRouter 和 AutoCodeHistoryRouter 引用 ✓

- [x] Step 2.4: 修改 `gin-vue-admin/server/initialize/router.go`
  - 移除 InitAutoCodeRouter 和 InitAutoCodeHistoryRouter 调用 ✓

### Phase 3: 前端文件删除

- [x] Step 3.1: 删除前端 AutoCode API 文件
  - `gin-vue-admin/web/src/api/autoCode.js` ✓

- [x] Step 3.2: 删除前端 systemTools 页面目录
  - `gin-vue-admin/web/src/view/systemTools/autoCode/` (整个目录) ✓
  - `gin-vue-admin/web/src/view/systemTools/autoCodeAdmin/` (整个目录) ✓
  - `gin-vue-admin/web/src/view/systemTools/autoPkg/` (整个目录) ✓
  - `gin-vue-admin/web/src/view/systemTools/autoPlug/` (整个目录) ✓
  - `gin-vue-admin/web/src/view/systemTools/formCreate/` (整个目录) ✓
  - `gin-vue-admin/web/src/view/systemTools/installPlugin/` (整个目录) ✓

### Phase 4: 废弃代码清理

- [x] Step 4.1: 删除 rm_file 目录
  - `gin-vue-admin/rm_file/` (整个目录) ✓

### Phase 5: 验证

- [x] Step 5.1: 验证后端编译
  - 在 `gin-vue-admin/server/` 目录执行 `go build` ✓
  - 额外修复: `initialize/ensure_tables.go` 和 `initialize/gorm.go` 中的 AutoCode 引用 ✓

- [x] Step 5.2: 验证前端构建
  - 在 `gin-vue-admin/web/` 目录执行 `npm run build` ✓
  - 额外修复: `chatTable.vue` 中移除了对已删除 `autoCode` API 的引用 ✓

---

## Summary
- **Total Steps**: 14
- **Delete Operations**: 10 steps
- **Modify Operations**: 4 steps
- **Verification**: 2 steps
- **Estimated Time**: 30-60 minutes
- **Status**: ✅ 全部完成
