<template>
  <div class="page">
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
        <el-form-item label="创建时间">
          <el-date-picker v-model="searchInfo.startCreatedAt" type="datetime" placeholder="开始时间"></el-date-picker>
          <el-date-picker v-model="searchInfo.endCreatedAt" type="datetime" placeholder="结束时间"></el-date-picker>
        </el-form-item>
        <el-form-item label="流量分类">
          <el-select v-model="searchInfo.category" style="width:194px">
            <el-option v-for="item in CategoryOption" :key="item.categoryName" :value="item.categoryName" :label="item.categoryName" />
          </el-select>
        </el-form-item>
        <el-form-item label="用户名">
          <el-select v-model="searchInfo.tag" clearable filterable style="width:194px">
            <el-option v-for="item in users" :key="item.ID" :value="item.ID" :label="item.nickName" />
          </el-select>
        </el-form-item>
        <el-form-item label="服务器Ip">
          <el-select v-model="searchInfo.server_ip" clearable filterable style="width:194px">
            <el-option v-for="item in srvs" :key="item.ip" :value="item.ip" :label="item.ip" />
          </el-select>
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
        <el-table-column align="left" label="标签" prop="tag" width="120" />
        <el-table-column align="left" label="用户" prop="username" width="120" />
        <el-table-column align="left" label="流量分类" prop="category" width="120" />
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
  deleteStat,
  findStat,
  getStatList,
  getStatCharts,
} from '@/api/stat'
import {
  getAllUserApi
} from '@/api/user'
import { getAllServerApi } from '@/api/server'
import { ElMessage, ElMessageBox } from 'element-plus'
import {ref, reactive, shallowRef, onMounted, nextTick, onUnmounted} from 'vue'
import EchartsLine from './statChart.vue'
import { setChartData } from "./common";
const formData = ref({
  category: '',
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
const searchInfo = ref({
  category: "user"
})

const onReset = () => {
  searchInfo.value = {
    category: "user"
  }
  getTableData()
  setChartData({...searchInfo.value})
}

// 搜索
const onSubmit = () => {
  page.value = 1
  pageSize.value = 10
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

const CategoryOption = ref([
  {
    categoryName: 'user'
  },
  {
    categoryName: 'inbound'
  },
  {
    categoryName: 'outbound'
  }
])

const users = ref([])
const getUsers = async() => {
  const res = await getAllUserApi()
  if (res.code === 0) {
    users.value = res.data.users
  }
}

const srvs = ref([])
const getSrvs = async() => {
  const res = await getAllServerApi()
  if (res.code === 0) {
    srvs.value = res.data.srvs
  }
}

const init = () => {
  getUsers()
  getSrvs()
}
init()

</script>
<style>
</style>