<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline" role="search" aria-label="服务器筛选" @keyup.enter="onSubmit">
      <el-form-item label="创建时间">
      <el-date-picker v-model="searchInfo.startCreatedAt" type="date" placeholder="开始时间" aria-label="服务器创建开始日期"></el-date-picker>
      <el-date-picker v-model="searchInfo.endCreatedAt" type="date" placeholder="结束时间" aria-label="服务器创建结束日期"></el-date-picker>
      </el-form-item>
        <el-form-item label="服务器ip">
         <el-input v-model="searchInfo.ip" placeholder="搜索条件" aria-label="服务器 IP 筛选" />

        </el-form-item>
        <el-form-item role="group" aria-label="服务器查询操作">
          <el-button type="primary" :icon="$gvaIcons.Search" @click="onSubmit">查询</el-button>
          <el-button :icon="$gvaIcons.Refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box" role="region" aria-label="服务器列表明细">
        <div class="gva-btn-list" role="group" aria-label="服务器列表操作">
            <el-button type="primary" :icon="$gvaIcons.Plus" @click="openDialog">新增</el-button>
            <el-popover v-model:visible="deleteVisible" placement="top" width="160">
            <p>确定要删除选中的服务器吗？</p>
            <div style="text-align: right; margin-top: 8px;">
                <el-button type="primary" link aria-label="取消批量删除服务器" @click="deleteVisible = false">取消</el-button>
                <el-button type="danger" aria-label="确认批量删除服务器" @click="onDelete">确定</el-button>
            </div>
            <template #reference>
                <el-button type="danger" :icon="$gvaIcons.Delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" aria-label="批量删除服务器" @click="deleteVisible = true">删除</el-button>
            </template>
            </el-popover>
        </div>
        <el-table
        ref="multipleTable"
        v-loading="tableLoading"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        :aria-busy="tableLoading"
        aria-label="服务器列表"
        :empty-text="tableLoading ? '加载中...' : '暂无数据'"
        row-key="ID"
        @selection-change="handleSelectionChange"
        >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="服务器ip" prop="ip" width="150" />
        <el-table-column align="left" label="端口" prop="port" width="60" />
        <el-table-column align="left" label="使用额度" width="100" >
            <template #default="scope">{{ formatFlow(scope.row.used_quota) }}</template>
        </el-table-column>
        <el-table-column align="left" label="总额度" width="100" >
            <template #default="scope">{{ scope.row.total_quota }}</template>
        </el-table-column>
        <el-table-column align="left" label="备注" prop="remark" width="120" />
        <el-table-column align="left" label="重置日期" prop="reset_date" width="100" />
        <el-table-column align="left" label="创建日期" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="操作">
            <template #default="scope">
                <div class="table-row-actions" role="group" aria-label="服务器行操作">
                    <el-button type="primary" link :icon="$gvaIcons.Document" class="table-button" @click="showServerConfig(scope.row)">查看配置</el-button>
                    <el-button type="primary" link :icon="$gvaIcons.Edit" class="table-button" @click="updateServerFunc(scope.row)">编辑</el-button>
                    <el-button type="warning" link :icon="$gvaIcons.Refresh" aria-label="重启服务器代理" @click="restartXray(scope.row)">代理重启</el-button>
                    <el-button type="danger" link :icon="$gvaIcons.Delete" aria-label="删除服务器" @click="deleteRow(scope.row)">删除</el-button>
                </div>
            </template>
        </el-table-column>
        </el-table>
        <div class="gva-pagination" role="navigation" aria-label="服务器列表分页">
            <el-pagination
              aria-label="服务器列表分页"
            layout="total, sizes, prev, pager, next, jumper"
            :current-page="page"
            :page-size="pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="total"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
            />
        </div>
    </div>
    <el-dialog v-model="dialogConfigVisible" :before-close="closeConfigDialog" title="配置信息">
      <el-form :model="configInfo" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
        <pre>{{ configInfo.content }}</pre>
      </el-form>
    </el-dialog>
    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" title="服务器信息">
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="100px">
        <el-form-item label="服务器ip:"  prop="ip" >
          <el-input v-model="formData.ip" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="端口:"  prop="port" >
          <el-input v-model.number="formData.port" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="备注:"  prop="remark" >
          <el-input v-model="formData.remark" :clearable="false"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="重置日期:"  prop="resetDate">
          <el-input v-model.number="formData.reset_date" :clearable="false"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="总额度(GB):"  prop="resetDate">
          <el-input v-model.number="formData.total_quota" :clearable="false"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="统计端口:"  prop="stat_port" >
          <el-input v-model.number="formData.stat_port" :clearable="true" placeholder="默认: 56611" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'Server'
}
</script>

<script setup>
import {
  createServer,
  deleteServer,
  deleteServerByIds,
  updateServer,
  getServerList,
  restartXrayApi
} from '@/api/server'

// 全量引入格式化工具 请按需保留
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import { buildServerPortChangeReminder } from './serverPortReminder.mjs'
import { devError } from '@/utils/devLogger'

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
        ip: '',
        port: 80,
        remark: '',
        total_quota: 1000,
        stat_port: 0,
        })
const originalFormData = ref(null)

// 验证规则
const rule  = reactive({
  ip: [
    { required: true, message: '请输入服务器IP', trigger: 'blur' },
  ],
  port: [
    { required: true, message: '请输入端口号', trigger: 'blur' },
  ],
  reset_date: [
    { required: true, message: '请输入流量重置时间', trigger: 'blur' },
  ],
})

const elFormRef = ref()


// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const tableLoading = ref(false)
const searchInfo = ref({})
const configInfo = ref({
  content: '',
})

// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// 搜索
const onSubmit = () => {
  page.value = 1
  pageSize.value = 10
  getTableData()
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// 修改页面容量
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async() => {
  tableLoading.value = true
  try {
    const table = await getServerList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
    if (table.code === 0) {
      tableData.value = table.data.list
      total.value = table.data.total
      page.value = table.data.page
      pageSize.value = table.data.pageSize
    }
  } finally {
    tableLoading.value = false
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () =>{
}

// 获取需要的字典 可能为空 按需保留
setOptions()

// 多选数据
const multipleSelection = ref([])
// 多选
const handleSelectionChange = (val) => {
    multipleSelection.value = val
}

// 删除行
const deleteRow = (row) => {
    ElMessageBox.confirm('确定要删除该服务器吗?', '删除服务器', {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
            deleteServerFunc(row)
        })
    }

const restartXray = (row) => {
  ElMessageBox.confirm('确定要重启该服务器代理吗?', '重启代理', {
    confirmButtonText: '重启',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    restartXrayFunc(row)
  })
}


// 批量删除控制标记
const deleteVisible = ref(false)

// 多选删除
const onDelete = async() => {
      const ids = []
      if (multipleSelection.value.length === 0) {
        ElMessage({
          type: 'warning',
          message: '请选择要删除的数据'
        })
        return
      }
      multipleSelection.value &&
        multipleSelection.value.map(item => {
          ids.push(item.ID)
        })
      const res = await deleteServerByIds({ ids })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功'
        })
        if (tableData.value.length === ids.length && page.value > 1) {
          page.value--
        }
        deleteVisible.value = false
        getTableData()
      }
    }

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

const dialogConfigVisible = ref(false)

const showServerConfig = async (row) => {
    configInfo.value.content = JSON.stringify(row.config, null, 2)
    dialogConfigVisible.value = true
}

const closeConfigDialog = () => {
  dialogConfigVisible.value = false
}

// 更新行
const updateServerFunc = async(row) => {
    try {
        const res = await updateServer({ ...row })
        if (res.code === 0) {
            // Only set type and show dialog if successful
            type.value = 'update'
            
            // Get the server data directly from the row parameter
            // This ensures we have all the necessary fields including total_quota
            formData.value = { ...row }
            originalFormData.value = { ...row }
            
            dialogFormVisible.value = true
        } else {
            // Handle API error response
            ElMessage({
                type: 'error',
                message: res.message || '更新失败'
            })
        }
    } catch (error) {
        // Handle unexpected errors
        devError('更新服务器时出错:', error)
        ElMessage({
            type: 'error',
            message: '更新服务器时发生错误'
        })
    }
}


// 删除行
const deleteServerFunc = async (row) => {
    const res = await deleteServer({ ID: row.ID })
    if (res.code === 0) {
        ElMessage({
                type: 'success',
                message: '删除成功'
            })
            if (tableData.value.length === 1 && page.value > 1) {
            page.value--
        }
        getTableData()
    }
}

const restartXrayFunc = async (row) => {
    const res = await restartXrayApi({ ID: row.ID, port: row.port, ip: row.ip })
    if (res?.code === 0) {
        ElMessage({
            type: 'success',
            message: '重启成功'
        })
    }
}

// 弹窗控制标记
const dialogFormVisible = ref(false)

// 打开弹窗
const openDialog = () => {
    type.value = 'create'
    originalFormData.value = null
    dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
    dialogFormVisible.value = false
    originalFormData.value = null
    formData.value = {
        ip: '',
        port: 0,
        remark: '',
        total_quota: 1000,
        stat_port: 0,
        }
}
const confirmPortChangeReminder = async () => {
  const reminder = buildServerPortChangeReminder(originalFormData.value, formData.value)
  if (!reminder) {
    return true
  }

  try {
    await ElMessageBox.confirm(reminder, '端口变更提醒', {
      confirmButtonText: '已同步，继续保存',
      cancelButtonText: '返回修改',
      type: 'warning'
    })
    return true
  } catch {
    return false
  }
}

// 弹窗确定
const enterDialog = async () => {
     elFormRef.value?.validate( async (valid) => {
             if (!valid) return
              if (!(await confirmPortChangeReminder())) return
              let res
              switch (type.value) {
                case 'create':
                  res = await createServer(formData.value)
                  break
                case 'update':
                  res = await updateServer(formData.value)
                  break
                default:
                  res = await createServer(formData.value)
                  break
              }
              if (res.code === 0) {
                ElMessage({
                  type: 'success',
                  message: '创建/更改成功'
                })
                closeDialog()
                getTableData()
              }
      })
}

const formatFlow = (value) => {
  const bytes = Number(value) || 0
  if (bytes >= 1024 * 1024 * 1024) { // 大于等于 1G 显示 G 后缀
    return (bytes / (1024 * 1024 * 1024)).toFixed(1) + ' G'
  } else if (bytes >= 1024 * 1024) { // 大于等于 1M 显示 M 后缀
    return (bytes / (1024 * 1024)).toFixed(1) + ' M'
  } else if (bytes >= 1024) { // 大于等于 1k 显示 k 后缀
    return (bytes / 1024).toFixed(1) + ' K'
  } else { // 小于 1k 直接返回原数值
    return bytes.toFixed(1)
  }
}

</script>

<style scoped>
/* 移动端优化 */
@media screen and (max-width: 768px) {
  /* 表格横向滚动优化 - 确保按钮固定 */
  .gva-table-box {
    position: relative;
    
    /* 按钮区域固定，不随表格滚动 */
    .gva-btn-list {
      position: relative;
      z-index: 2;
      background: #fff;
      padding: 0.5rem;
      border-bottom: 1px solid #ebeef5;
    }
    
    /* 表格容器可滚动 */
    .el-table {
      overflow-x: auto;
      -webkit-overflow-scrolling: touch;
    }
    
    /* 添加渐变提示 */
    &::before {
      content: '';
      position: absolute;
      right: 0;
      top: 4rem; /* 调整位置，避免覆盖按钮 */
      bottom: 3rem;
      width: 2rem;
      background: linear-gradient(to left, rgba(255,255,255,0.9), transparent);
      pointer-events: none;
      z-index: 1;
    }
  }
  
  /* 配置信息弹窗 */
  .el-dialog {
    pre {
      font-size: 0.75rem;
      overflow-x: auto;
      white-space: pre-wrap;
      word-wrap: break-word;
    }
  }
}
</style>
