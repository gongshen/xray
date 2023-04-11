<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
      <el-form-item label="创建时间">
      <el-date-picker v-model="searchInfo.startCreatedAt" type="datetime" placeholder="开始时间"></el-date-picker>
      <el-date-picker v-model="searchInfo.endCreatedAt" type="datetime" placeholder="结束时间"></el-date-picker>
      </el-form-item>
        <el-form-item label="服务器ip">
         <el-input v-model="searchInfo.ip" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
        <div class="gva-btn-list">
            <el-button type="primary" icon="plus" @click="openDialog">新增</el-button>
            <el-popover v-model:visible="deleteVisible" placement="top" width="160">
            <p>确定要删除吗？</p>
            <div style="text-align: right; margin-top: 8px;">
                <el-button type="primary" link @click="deleteVisible = false">取消</el-button>
                <el-button type="primary" @click="onDelete">确定</el-button>
            </div>
            <template #reference>
                <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="deleteVisible = true">删除</el-button>
            </template>
            </el-popover>
        </div>
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
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
            <template #default="scope">{{ formatFlow(scope.row.total_quota) }}</template>
        </el-table-column>
        <el-table-column align="left" label="备注" prop="remark" width="120" />
        <el-table-column align="left" label="重置日期" prop="reset_date" width="100" />
        <el-table-column align="left" label="创建日期" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="按钮组">
            <template #default="scope">
                <el-button type="primary" link icon="document" class="table-button" @click="showServerConfig(scope.row)">查看配置</el-button>
<!--            <el-button type="primary" link icon="edit" class="table-button" @click="updateServerFunc(scope.row)">变更</el-button>-->
            <el-button type="primary" link icon="refresh" @click="restartXray(scope.row)">代理重启</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
            </template>
        </el-table-column>
        </el-table>
        <div class="gva-pagination">
            <el-pagination
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
    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" title="弹窗操作">
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
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
  findServer,
  getServerList,
  restartXrayApi
} from '@/api/server'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
        ip: '',
        port: 80,
        remark: '',
        total_quota: 1000,
        })

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
  const table = await getServerList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
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
    ElMessageBox.confirm('确定要删除吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
            deleteServerFunc(row)
        })
    }

const restartXray = (row) => {
  ElMessageBox.confirm('确定要重启代理吗?', '提示', {
    confirmButtonText: '确定',
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
    const res = await findServer({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data.reserver
        dialogFormVisible.value = true
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
    if (res.code === 0) {
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
    dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
        ip: '',
        port: 0,
        remark: '',
        }
}
// 弹窗确定
const enterDialog = async () => {
     elFormRef.value?.validate( async (valid) => {
             if (!valid) return
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
  if (value >= 1024 * 1024 * 1024) { // 大于等于 1G 显示 G 后缀
    return (value / (1024 * 1024 * 1024)).toFixed(1) + ' G'
  } else if (value >= 1024 * 1024) { // 大于等于 1M 显示 M 后缀
    return (value / (1024 * 1024)).toFixed(1) + ' M'
  } else if (value >= 1024) { // 大于等于 1k 显示 k 后缀
    return (value / 1024).toFixed(1) + ' K'
  } else { // 小于 1k 直接返回原数值
    return value.toFixed(1)
  }
}

</script>

<style>
</style>
