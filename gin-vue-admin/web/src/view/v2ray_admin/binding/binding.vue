<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
      <el-form-item label="创建时间">
      <el-date-picker v-model="searchInfo.startCreatedAt" type="date" placeholder="开始时间"></el-date-picker>
       —
      <el-date-picker v-model="searchInfo.endCreatedAt" type="date" placeholder="结束时间"></el-date-picker>
      </el-form-item>
        <el-form-item label="服务器ip">
         <el-select v-model="searchInfo.server_id" clearable filterable style="width:194px">
           <el-option v-for="item in srvs" :key="item.ID" :value="item.ID" :label="item.ip" />
         </el-select>
         <el-icon class="auto-icon" @click="getSrvs"><refresh /></el-icon>
        </el-form-item>
        <el-form-item label="用户名">
         <el-select v-model="searchInfo.user_id" clearable filterable style="width:194px">
           <el-option v-for="item in users" :key="item.ID" :value="item.ID" :label="item.nickName" />
         </el-select>
          <el-icon class="auto-icon" @click="getUsers"><refresh /></el-icon>
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
        <el-table-column align="left" label="用户名" width="120">
            <template #default="scope">{{ scope.row.user.nickName }}</template>
        </el-table-column>
        <el-table-column align="left" label="服务器ip" width="140">
            <template #default="scope">{{ scope.row.server.ip }}</template>
        </el-table-column>
        <el-table-column align="left" label="日期" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="是否限流" width="100">
            <template #default="scope">{{ scope.row.is_limited ? '是' : '否' }}</template>
        </el-table-column>
        <el-table-column align="left" label="按钮组">
            <template #default="scope">
            <el-button type="primary" link icon="share" class="table-button" @click="shareBindingFunc(scope.row)">分享</el-button>
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
    <el-dialog v-model="shareFormVisible" :before-close="closeShareDialog" title="弹窗操作">
      <el-form :model="shareInfo" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
        <h2>shadowrocket,qv2ray,v2rayXS</h2>
        <div>
          <img :src="shareInfo.share1" alt="二维码" style="width: 300px; height: 300px;"/>
        </div>
        <el-button class="btn" :data-clipboard-text="shareInfo.share1_link">复制</el-button>
        <el-divider></el-divider>
        <h2>v2rayN,v2rayNG,v2rayXS</h2>
        <div>
          <img :src="shareInfo.share2" alt="二维码" style="width: 300px; height: 300px;"/>
        </div>
        <el-button class="btn" :data-clipboard-text="shareInfo.share2_link">复制</el-button>
      </el-form>
    </el-dialog>
    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" title="弹窗操作">
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
        <el-form-item label="服务器ip:">
          <el-select v-model="formData.server_id" clearable filterable style="width:194px">
            <el-option v-for="item in srvs" :key="item.ID" :value="item.ID" :label="item.ip" />
          </el-select>
        </el-form-item>
        <el-form-item label="用户名:">
          <el-select v-model="formData.user_id" clearable filterable style="width:194px">
            <el-option v-for="item in users" :key="item.ID" :value="item.ID" :label="item.nickName" />
          </el-select>
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
  name: 'Binding'
}
</script>

<script setup>
import {
  createBinding,
  deleteBinding,
  deleteBindingByIds,
  updateBinding,
  findBinding,
  getBindingList,
  shareBinding
} from '@/api/binding'
import { getAllServerApi } from '@/api/server'
import { getAllUserApi } from '@/api/user'
import  QRCode  from 'qrcode'
import ClipboardJS from 'clipboard';

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import {ref, reactive, onMounted} from 'vue'

const clipboard = new ClipboardJS('.btn');

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
        server_id: '',
        user_id: '',
        alter_id: 64,
        })

// 验证规则
const rule = reactive({})

const elFormRef = ref()


// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const shareInfo = ref({
  share1: '',
  share1_link: '',
  share2: '',
  share2_link: '',
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
  const table = await getBindingList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
            deleteBindingFunc(row)
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
      const res = await deleteBindingByIds({ ids })
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

// 分享标记
const shareFormVisible = ref(false)

const shareBindingFunc = async(row) => {
    const res = await shareBinding({ ID: row.ID })
    if (res.code === 0) {
      console.log("aaa:", res.data.share1)
        shareInfo.value.share1_link = res.data.share1
        shareInfo.value.share2_link = res.data.share2
        QRCode.toDataURL(res.data.share1)
            .then((url) => {
                shareInfo.value.share1 = url
            })
        QRCode.toDataURL(res.data.share2)
            .then((url) => {
                shareInfo.value.share2 = url
            })
        shareFormVisible.value = true
    }
}

const closeShareDialog = () => {
  shareFormVisible.value = false
}

// 删除行
const deleteBindingFunc = async (row) => {
    const res = await deleteBinding({ ID: row.ID })
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
        server_id: '',
        user_id: '',
        }
}
// 弹窗确定
const enterDialog = async () => {
     elFormRef.value?.validate( async (valid) => {
             if (!valid) return
              let res
              switch (type.value) {
                case 'create':
                  res = await createBinding(formData.value)
                  break
                case 'update':
                  res = await updateBinding(formData.value)
                  break
                default:
                  res = await createBinding(formData.value)
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

const srvs = ref([])
const getSrvs = async() => {
  const res = await getAllServerApi()
  if (res.code === 0) {
    srvs.value = res.data.srvs
  }
}

const users = ref([])
const getUsers = async() => {
  const res = await getAllUserApi()
  if (res.code === 0) {
    users.value = res.data.users
  }
}

const init = () => {
  getSrvs()
  getUsers()
}
init()

</script>

<style>
</style>
