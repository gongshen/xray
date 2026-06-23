<template>
  <div class="page">
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
        <el-form-item label="创建时间">
          <el-date-picker v-model="searchInfo.startCreatedAt" type="date" placeholder="开始时间" :default-value="monthAgo"></el-date-picker>
          <el-date-picker v-model="searchInfo.endCreatedAt" type="date" placeholder="结束时间" :default-value="today"></el-date-picker>
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

    <!-- 图表区域 -->
    <div class="charts-section">
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
            <div ref="echart" class="chart-container trend-chart"></div>
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
            <div ref="rank_echart" class="chart-container rank-chart"></div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 详细数据表格 -->
    <div class="gva-table-box">
      <div class="table-header">
        <h3>详细流量记录</h3>
        <div class="table-summary">
          <el-tag type="primary" effect="light">总流量 {{ formatFlow(chartData.total) }}</el-tag>
          <el-tag type="info" effect="plain">共 {{ total }} 条记录</el-tag>
        </div>
      </div>
      
      <el-table
          ref="multipleTable"
          v-loading="tableLoading"
          style="width: 100%; min-height: 200px;"
          tooltip-effect="dark"
          :data="tableData"
          row-key="ID"
          stripe
          :empty-text="tableData.length === 0 ? '暂无数据' : ''"
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
import { TrendCharts, Trophy } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { useChartData, setChartData } from "./common"
import { bindWindowEvent } from '@/utils/eventLifecycle.mjs'
import {
  formatFlow,
  getDateRangeText as formatDateRangeText,
  getTrafficTagType,
  normalizeDateOnlyToUtcIso,
} from './statTraffic.mjs'
import {
  buildRankChartOptions,
  buildTrendChartOptions,
} from './statChartOptions.mjs'

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const tableLoading = ref(false)
const today = new Date()
const monthAgo = new Date(today.getTime() - 30 * 24 * 60 * 60 * 1000) // 1个月前
const searchInfo = ref({
  startCreatedAt: monthAgo.toISOString(),
  endCreatedAt: today.toISOString()
})

// 图表相关
const chart = shallowRef(null)
const rankChart = shallowRef(null)
let disposeResize = null
const echart = ref(null)
const rank_echart = ref(null)
const chartData = useChartData()
const dateRangeText = computed(() => formatDateRangeText(searchInfo.value))

const onReset = () => {
  // 重置为近1个月
  const today = new Date()
  const monthAgo = new Date(today.getTime() - 30 * 24 * 60 * 60 * 1000)
  searchInfo.value = {
    startCreatedAt: monthAgo.toISOString(),
    endCreatedAt: today.toISOString()
  }
  getTableData()
  setChartData({...searchInfo.value})
}

// 搜索
const onSubmit = () => {
  page.value = 1
  pageSize.value = 10
  // 将搜索时间转换为 UTC 时间
  if (searchInfo.value.startCreatedAt) {
    searchInfo.value.startCreatedAt = normalizeDateOnlyToUtcIso(searchInfo.value.startCreatedAt)
  }
  if (searchInfo.value.endCreatedAt) {
    searchInfo.value.endCreatedAt = normalizeDateOnlyToUtcIso(searchInfo.value.endCreatedAt)
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
  tableLoading.value = true
  try {
    const table = await getStatList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })

    if (table.code === 0) {
      // 确保数据结构正确
      const list = table.data?.list || []
      tableData.value = list
      total.value = table.data?.total || 0
      page.value = table.data?.page || 1
      pageSize.value = table.data?.pageSize || 10

      // 如果有数据但表格不显示，可能是响应式问题
      if (list.length > 0) {
        // 强制触发响应式更新
        tableData.value = [...list]
      }
    } else {
      console.error('getStatList error:', table.msg)
      ElMessage.error(table.msg || '获取数据失败')
      tableData.value = []
      total.value = 0
    }
  } catch (error) {
    console.error('getTableData error:', error)
    ElMessage.error('网络请求失败')
    tableData.value = []
    total.value = 0
  } finally {
    tableLoading.value = false
  }
}

// 初始化图表
const initChart = () => {
  if (!echart.value || !rank_echart.value) {
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
  initChart()
  disposeResize = bindWindowEvent(window, 'resize', handleResize)
})

onUnmounted(() => {
  disposeResize?.()
  disposeResize = null
  disposeCharts()
})
</script>

<style scoped>
.page {
  padding: 20px;
  background: var(--gva-color-page-bg);
  min-height: 100vh;
}

/* 图表区域 */
.charts-section {
  margin-bottom: 24px;
}

.chart-card {
  border-radius: var(--gva-radius-panel);
  overflow: hidden;
  border: none;
  box-shadow: var(--gva-shadow-panel);
  transition: all 0.3s ease;
}

.chart-card:hover {
  box-shadow: var(--gva-shadow-panel-hover);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0;
}

.card-title {
  display: flex;
  align-items: center;
  font-size: 16px;
  font-weight: 600;
  color: var(--gva-color-text-strong);
}

.card-title .el-icon {
  margin-right: 8px;
  font-size: 18px;
  color: var(--gva-color-brand-secondary);
}

.chart-container {
  height: 360px;
  width: 100%;
}

.trend-chart {
  background: linear-gradient(180deg, var(--gva-color-chart-trend-bg) 0%, var(--gva-color-panel-bg) 100%);
}

.rank-chart {
  background: linear-gradient(180deg, var(--gva-color-chart-rank-bg) 0%, var(--gva-color-panel-bg) 100%);
}

/* 表格区域 */
.gva-table-box {
  background: var(--gva-color-panel-bg);
  border-radius: var(--gva-radius-panel);
  padding: 24px;
  box-shadow: var(--gva-shadow-panel);
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--gva-color-border-muted);
}

.table-header h3 {
  margin: 0;
  color: var(--gva-color-text-strong);
  font-size: 18px;
  font-weight: 600;
}

.table-summary {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  justify-content: flex-end;
}

/* 表格单元格样式 */
.traffic-down {
  color: var(--gva-color-traffic-down);
  font-weight: 500;
}

.traffic-up {
  color: var(--gva-color-traffic-up);
  font-weight: 500;
}

.date-cell {
  color: var(--gva-color-text-muted);
  font-size: 13px;
}

/* 移动端优化 */
@media screen and (max-width: 768px) {
  .page {
    padding: 10px;
  }

  .charts-section {
    margin-bottom: 16px;
  }

  .chart-container {
    height: 280px;
  }

  .gva-table-box {
    padding: 16px;
    border-radius: var(--gva-radius-panel);
  }

  .table-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .table-header h3 {
    font-size: 16px;
  }

  .table-summary {
    justify-content: flex-start;
  }

  /* 表格横向滚动优化 */
  .gva-table-box {
    position: relative;
    
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
      top: 80px;
      bottom: 60px;
      width: 2rem;
      background: linear-gradient(to left, rgba(255, 255, 255, 0.9), transparent);
      pointer-events: none;
      z-index: 1;
    }
  }

}

/* 响应式图表 */
@media screen and (max-width: 1200px) {
  .charts-section .el-col {
    margin-bottom: 16px;
  }
}

/* 深色模式适配 */
@media (prefers-color-scheme: dark) {
  .page {
    background: var(--gva-color-dark-page-bg);
  }

  .chart-card,
  .gva-table-box {
    background: var(--gva-color-dark-panel-bg);
    color: var(--gva-color-dark-text);
  }

  .card-title {
    color: var(--gva-color-dark-text);
  }

  .table-header h3 {
    color: var(--gva-color-dark-text);
  }

  .table-header {
    border-bottom-color: var(--gva-color-dark-border);
  }
}

/* 动画效果 */
.chart-card,
.gva-table-box {
  animation: fadeInUp 0.6s ease-out;
}

@media (prefers-reduced-motion: reduce) {
  .chart-card,
  .gva-table-box {
    animation: none;
  }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 表格斑马纹优化 */
:deep(.el-table--striped .el-table__body tr.el-table__row--striped td) {
  background: var(--gva-color-panel-muted-bg);
}

/* 分页样式优化 */
.gva-pagination {
  margin-top: 24px;
  display: flex;
  justify-content: center;
}

:deep(.el-pagination) {
  --el-pagination-button-color: var(--gva-color-text-regular);
  --el-pagination-hover-color: var(--gva-color-brand-secondary);
}
</style>
