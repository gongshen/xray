<template>
  <div class="page">
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
        <el-form-item label="创建时间">
          <el-date-picker v-model="searchInfo.startCreatedAt" type="date" placeholder="开始时间"></el-date-picker>
          <el-date-picker v-model="searchInfo.endCreatedAt" type="date" placeholder="结束时间"></el-date-picker>
        </el-form-item>
        <el-form-item>
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
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="用户" prop="username" width="120" />
        <el-table-column align="left" label="服务器Ip" prop="server_ip" width="200" />
        <el-table-column align="left" label="下行流量" prop="down" width="120" />
        <el-table-column align="left" label="上行流量" prop="up" width="120" />
        <el-table-column align="left" label="总流量" prop="total" width="120" />
        <el-table-column align="left" label="日期" prop="created_at" width="180"/>
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
    <div class="gva-card-box">
      <el-row>
        <el-col>
          <echarts-line />
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script>

export default {
  name: 'Stat'
}
</script>

<script setup>
import {
  getStatList,
} from '@/api/v2ray_stat'
import {ref, reactive,} from 'vue'
import EchartsLine from './statChart.vue'
import { setChartData } from "./common";
const formData = ref({
  tag: '',
  down: '',
  up: '',
  total: '',
})
const rule = reactive({
  category : [{
    required: true,
    message: '',
    trigger: ['input','blur'],
  }],
})
const elFormRef = ref()
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const today = new Date().toISOString()
const searchInfo = ref({
  startCreatedAt: today,
  endCreatedAt: today
})

const onReset = () => {
  getTableData()
  setChartData({...searchInfo.value})
}

// 搜索
const onSubmit = () => {
  page.value = 1
  pageSize.value = 10
  if (searchInfo.value.startCreatedAt) {
    const startDate = new Date(searchInfo.value.startCreatedAt)
    const utcStartDate = new Date(Date.UTC(startDate.getFullYear(), startDate.getMonth(), startDate.getDate()))
    searchInfo.value.startCreatedAt = utcStartDate.toISOString()
  }
  if (searchInfo.value.endCreatedAt) {
    const endDate = new Date(searchInfo.value.endCreatedAt)
    const utcEndDate = new Date(Date.UTC(endDate.getFullYear(), endDate.getMonth(), endDate.getDate()))
    searchInfo.value.endCreatedAt = utcEndDate.toISOString()
  }
  getTableData()
  setChartData({...searchInfo.value})
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
  const table = await getStatList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()
setChartData({...searchInfo.value})

const users = ref([])

</script>
<style>
</style>