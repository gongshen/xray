<template>
  <div class="page">
    <div class="gva-search-box" role="search" aria-label="流量统计筛选">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
        <el-form-item label="创建时间">
          <el-date-picker v-model="searchInfo.startCreatedAt" type="date" placeholder="开始时间" aria-label="统计开始日期" :default-value="monthAgo"></el-date-picker>
          <el-date-picker v-model="searchInfo.endCreatedAt" type="date" placeholder="结束时间" aria-label="统计结束日期" :default-value="today"></el-date-picker>
        </el-form-item>
        <el-form-item label="用户名">
          <el-select v-model="searchInfo.tag" clearable filterable class="stat-filter-select" aria-label="统计用户筛选">
            <el-option v-for="item in users" :key="item.ID" :value="item.ID" :label="item.nickName" />
          </el-select>
        </el-form-item>
        <el-form-item label="服务器Ip">
          <el-select v-model="searchInfo.server_ip" clearable filterable class="stat-filter-select" aria-label="统计服务器筛选">
            <el-option v-for="item in srvs" :key="item.ip" :value="item.ip" :label="item.ip" />
          </el-select>
        </el-form-item>
        <el-form-item class="stat-action-group" role="group" aria-label="统计查询操作">
          <el-button type="primary" :icon="Search" @click="onSubmit">查询</el-button>
          <el-button :icon="Refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 图表区域 -->
    <div class="charts-section" role="region" aria-label="流量统计图表">
      <el-row :gutter="20">
        <el-col :xs="24" :lg="14">
          <el-card class="chart-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <span class="card-title">
                  <el-icon><TrendCharts /></el-icon>
                  流量趋势
                </span>
                <el-tag type="info" size="small">{{ dateRangeText }}</el-tag>
              </div>
            </template>
            <div ref="echart" class="chart-container trend-chart" role="img" :aria-label="trafficTrendChartLabel"></div>
          </el-card>
        </el-col>
        <el-col :xs="24" :lg="10">
          <el-card class="chart-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <span class="card-title">
                  <el-icon><Trophy /></el-icon>
                  流量排行榜
                </span>
                <el-tag type="success" size="small">TOP 10</el-tag>
              </div>
            </template>
            <div ref="rank_echart" class="chart-container rank-chart" role="img" :aria-label="trafficRankChartLabel"></div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 详细数据表格 -->
    <div class="gva-table-box" role="region" aria-label="流量记录明细">
      <div class="table-header">
        <h3>详细流量记录</h3>
        <div class="table-summary" role="status" aria-live="polite" aria-label="流量统计摘要">
          <el-tag type="primary" effect="light">总流量 {{ formatFlow(chartData.total) }}</el-tag>
          <el-tag type="info" effect="plain">共 {{ total }} 条记录</el-tag>
        </div>
      </div>
      
      <el-table
          ref="multipleTable"
          class="stat-table"
          v-loading="tableLoading"
          :aria-busy="tableLoading"
          aria-label="详细流量记录表"
          tooltip-effect="dark"
          :data="tableData"
          row-key="ID"
          stripe
          :empty-text="tableLoading ? '加载中...' : (tableData.length === 0 ? '暂无数据' : '')"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="用户" prop="username" width="120">
          <template #default="scope">
            <strong>{{ scope.row.username }}</strong>
          </template>
        </el-table-column>
        <el-table-column align="left" label="服务器" prop="server_ip" width="200">
          <template #default="scope">
            <el-tag type="info" size="small">{{ scope.row.server_ip }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="下行流量" prop="down" width="120">
          <template #default="scope">
            <span class="traffic-down">{{ scope.row.down }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="上行流量" prop="up" width="120">
          <template #default="scope">
            <span class="traffic-up">{{ scope.row.up }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="总流量" width="120">
          <template #default="scope">
            <el-tag :type="getTrafficTagType(scope.row.total)" size="small">
              {{ scope.row.total }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="日期" prop="created_at" width="180">
          <template #default="scope">
            <span class="date-cell">{{ scope.row.created_at }}</span>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination" role="navigation" aria-label="流量记录分页">
        <el-pagination
          aria-label="统计列表分页"
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
} from '@/api/stat'
import {
  getAllUserApi
} from '@/api/user'
import { getAllServerApi } from '@/api/server'
import { ElMessage } from 'element-plus'
import { computed, ref, shallowRef, onMounted, nextTick, onUnmounted, watch } from 'vue'
import { Refresh, Search, TrendCharts, Trophy } from '@element-plus/icons-vue'
import { useChartData, setChartData } from "./common"
import { bindWindowEvent } from '@/utils/eventLifecycle.mjs'
import { devError } from '@/utils/devLogger'
import { loadEcharts } from '@/utils/loadEcharts'
import {
  createDefaultTrafficSearchRange,
  formatFlow,
  getDateRangeText as formatDateRangeText,
  getTrafficTagType,
  normalizeTrafficSearchRange,
} from './statTraffic.mjs'
import {
  buildRankChartOptions,
  buildTrendChartOptions,
} from './statChartOptions.mjs'
import { normalizeStatTableResponse } from './statTableData.mjs'

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const tableLoading = ref(false)
const initialSearchRange = createDefaultTrafficSearchRange()
const today = new Date(initialSearchRange.endCreatedAt)
const monthAgo = new Date(initialSearchRange.startCreatedAt)
const searchInfo = ref({ ...initialSearchRange })

// 图表相关
const chart = shallowRef(null)
const rankChart = shallowRef(null)
let disposeResize = null
let isUnmounted = false
let tableRequestId = 0
let usersRequestId = 0
let srvsRequestId = 0
const echart = ref(null)
const rank_echart = ref(null)
const chartData = useChartData()
const dateRangeText = computed(() => formatDateRangeText(searchInfo.value))
const trafficTrendChartLabel = computed(() => `流量趋势图表，${dateRangeText.value}`)
const trafficRankChartLabel = computed(() => `流量排行榜图表，${dateRangeText.value}`)

const onReset = () => {
  // 重置为近1个月
  searchInfo.value = createDefaultTrafficSearchRange()
  getTableData()
  setChartData({...searchInfo.value})
}

// 搜索
const onSubmit = () => {
  page.value = 1
  pageSize.value = 10
  // 将搜索时间转换为 UTC 时间
  searchInfo.value = normalizeTrafficSearchRange(searchInfo.value)
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
  const requestId = ++tableRequestId
  tableLoading.value = true
  try {
    const table = await getStatList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
    if (requestId !== tableRequestId || isUnmounted) {
      return
    }
    const tableState = normalizeStatTableResponse(table, { page: page.value, pageSize: pageSize.value })

    if (tableState.ok) {
      tableData.value = tableState.list
      total.value = tableState.total
      page.value = tableState.page
      pageSize.value = tableState.pageSize
    } else {
      devError('getStatList error:', table.msg)
      ElMessage.error(tableState.message || '\u83b7\u53d6\u6570\u636e\u5931\u8d25')
      tableData.value = tableState.list
      total.value = tableState.total
    }
  } catch (error) {
    if (requestId === tableRequestId && !isUnmounted) {
      devError('getTableData error:', error)
      ElMessage.error('\u7f51\u7edc\u8bf7\u6c42\u5931\u8d25')
      tableData.value = []
      total.value = 0
    }
  } finally {
    if (requestId === tableRequestId && !isUnmounted) {
      tableLoading.value = false
    }
  }
}

// 初始化图表
const initChart = async() => {
  if (!echart.value || !rank_echart.value) {
    return
  }

  const echarts = await loadEcharts()
  if (!echart.value || !rank_echart.value || isUnmounted) {
    return
  }

  chart.value = echarts.init(echart.value)
  rankChart.value = echarts.init(rank_echart.value)
  setOptions(chartData)
}

// 设置图表选项
const setOptions = (data) => {
  chart.value?.setOption(buildTrendChartOptions(data))
  rankChart.value?.setOption(buildRankChartOptions(data))
}

const users = ref([])

const getUsers = async() => {
  const requestId = ++usersRequestId
  const res = await getAllUserApi()
  if (requestId === usersRequestId && !isUnmounted && res.code === 0) {
    users.value = res.data.users || []
  }
}

const srvs = ref([])

const getSrvs = async() => {
  const requestId = ++srvsRequestId
  const res = await getAllServerApi()
  if (requestId === srvsRequestId && !isUnmounted && res.code === 0) {
    srvs.value = res.data.srvs || []
  }
}

const init = () => {
  getUsers()
  getSrvs()
}

// 监听图表数据变化
watch(() => chartData, (newData) => {
  if (chart.value && rankChart.value) {
    setOptions(newData)
  }
}, {
  deep: true
})

const handleResize = () => {
  chart.value?.resize()
  rankChart.value?.resize()
}

const disposeCharts = () => {
  if (chart.value) {
    chart.value.dispose()
    chart.value = null
  }
  if (rankChart.value) {
    rankChart.value.dispose()
    rankChart.value = null
  }
}

onMounted(async () => {
  init()
  getTableData()
  setChartData({...searchInfo.value})
  
  await nextTick()
  await initChart()
  if (!isUnmounted) {
    disposeResize = bindWindowEvent(window, 'resize', handleResize)
  }
})

onUnmounted(() => {
  isUnmounted = true
  tableRequestId++
  usersRequestId++
  srvsRequestId++
  disposeResize?.()
  disposeResize = null
  disposeCharts()
})
</script>

<style scoped lang="scss" src="./statPage.scss"></style>
