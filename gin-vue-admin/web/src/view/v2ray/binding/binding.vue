<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
      <el-form-item label="创建时间">
      <el-date-picker v-model="searchInfo.startCreatedAt" type="datetime" placeholder="开始时间"></el-date-picker>
      <el-date-picker v-model="searchInfo.endCreatedAt" type="datetime" placeholder="结束时间"></el-date-picker>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
      </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
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
          <el-table-column align="left" label="按钮组">
            <template #default="scope">
            <el-button type="primary" link icon="share" class="table-button" @click="shareBindingFunc(scope.row)">分享</el-button>
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
  </div>
</template>

<script>
export default {
  name: 'Binding'
}
</script>

<script setup>
import {
  getBindingList,
  shareBinding
} from '@/api/v2ray_binding'
import  QRCode  from 'qrcode'
import ClipboardJS from 'clipboard';

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

const clipboard = new ClipboardJS('.btn');

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

</script>

<style>
</style>
