<template>
  <div class="page">
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
        <el-form-item label="创建时间">
          <el-date-picker v-model="searchInfo.startCreatedAt" type="date" placeholder="开始时间" :default-value="sevenDaysAgo"></el-date-picker>
          <el-date-picker v-model="searchInfo.endCreatedAt" type="date" placeholder="结束时间" :default-value="today"></el-date-picker>
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
        <el-col :xs="24" :lg="24">
          <el-card class="chart-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <span class="card-title">
                  <el-icon><TrendCharts /></el-icon>
                  流量趋势
                </span>
                <el-tag type="info" size="small">{{ getDateRangeText() }}</el-tag>
              </div>
            </template>
            <div ref="echart" class="chart-container trend-chart"></div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 详细数据表格 -->
    <div class="gva-table-box">
      <div class="table-header">
        <h3>详细流量记录</h3>
        <el-tag>共 {{ total }} 条记录</el-tag>
      </div>
      <el-table
          ref="multipleTable"
          style="width: 100%"
          tooltip-effect="dark"
          :data="tableData"
          row-key="ID"
          stripe
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="用户" prop="username" width="120">
          <template #default="scope">
            <div class="user-cell">
              <el-avatar :size="24" class="user-avatar">
                {{ scope.row.username?.charAt(0)?.toUpperCase() || 'U' }}
              </el-avatar>
              <span>{{ scope.row.username }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="服务器" prop="server_ip" width="200">
          <template #default="scope">
            <el-tag type="info" size="small">{{ scope.row.server_ip }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="下行流量" prop="down" width="120">
          <template #default="scope">
            <span class="traffic-down">{{ formatFlow(scope.row.down) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="上行流量" prop="up" width="120">
          <template #default="scope">
            <span class="traffic-up">{{ formatFlow(scope.row.up) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="总流量" width="120">
          <template #default="scope">
            <el-tag :type="getTrafficTagType(scope.row.down + scope.row.up)" size="small">
              {{ formatFlow(scope.row.down + scope.row.up) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="日期" prop="created_at" width="180">
          <template #default="scope">
            <span class="date-cell">{{ formatDate(scope.row.created_at) }}</span>
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
} from '@/api/v2ray_stat'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, shallowRef, onMounted, nextTick, onUnmounted, computed, watch } from 'vue'
import { TrendCharts } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { useChartData, setChartData } from "./common"

const formData = ref({
  tag: '',
  down: '',
  up: '',
  total: '',
})

const elFormRef = ref()
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const today = new Date()
const sevenDaysAgo = new Date(today.getTime() - 6 * 24 * 60 * 60 * 1000) // 7天前（包含今天共7天）
const searchInfo = ref({
  startCreatedAt: sevenDaysAgo.toISOString(),
  endCreatedAt: today.toISOString()
})

// 图表相关
const chart = shallowRef(null)
const echart = ref(null)
const chartData = useChartData()

// 格式化流量
const formatFlow = (value) => {
  if (!value || value === 0) return '0 B'
  if (value >= 1024 * 1024 * 1024) {
    return (value / (1024 * 1024 * 1024)).toFixed(1) + ' GB'
  } else if (value >= 1024 * 1024) {
    return (value / (1024 * 1024)).toFixed(1) + ' MB'
  } else if (value >= 1024) {
    return (value / 1024).toFixed(1) + ' KB'
  } else {
    return value.toFixed(1) + ' B'
  }
}

// 格式化日期
const formatDate = (timestamp) => {
  if (!timestamp) return '-'
  // 如果是秒级时间戳，转换为毫秒
  const date = new Date(timestamp * 1000)
  return date.toLocaleDateString('zh-CN') + ' ' + date.toLocaleTimeString('zh-CN', { hour12: false })
}

// 获取流量标签类型
const getTrafficTagType = (traffic) => {
  if (traffic >= 1024 * 1024 * 1024) return 'danger'  // GB级别
  if (traffic >= 100 * 1024 * 1024) return 'warning'  // 100MB以上
  if (traffic >= 10 * 1024 * 1024) return 'success'   // 10MB以上
  return 'info'  // 其他
}

// 获取日期范围文本
const getDateRangeText = () => {
  if (!searchInfo.value.startCreatedAt || !searchInfo.value.endCreatedAt) {
    return '近7天'
  }
  const start = new Date(searchInfo.value.startCreatedAt).toLocaleDateString('zh-CN')
  const end = new Date(searchInfo.value.endCreatedAt).toLocaleDateString('zh-CN')
  return start === end ? start : `${start} - ${end}`
}

const onReset = () => {
  // 重置为近7天
  const today = new Date()
  const sevenDaysAgo = new Date(today.getTime() - 6 * 24 * 60 * 60 * 1000)
  searchInfo.value = {
    startCreatedAt: sevenDaysAgo.toISOString(),
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
  try {
    const table = await getStatList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
    if (table.code === 0) {
      // 确保数据结构正确
      tableData.value = table.data?.list || []
      total.value = table.data?.total || 0
      page.value = table.data?.page || 1
      pageSize.value = table.data?.pageSize || 10
    } else {
      console.error('getStatList error:', table.msg)
      tableData.value = []
      total.value = 0
    }
  } catch (error) {
    console.error('getTableData error:', error)
    tableData.value = []
    total.value = 0
  }
}

// 初始化图表
const initChart = () => {
  chart.value = echarts.init(echart.value)
  setOptions(chartData)
}

// 设置图表选项
const setOptions = (data) => {
  // 流量趋势图
  chart.value?.setOption({
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
        label: {
          backgroundColor: '#6a7985'
        }
      },
      formatter: (params) => {
        let result = params[0].name + '<br>'
        params.forEach(item => {
          result += item.marker + ' ' + item.seriesName + ' : ' + formatFlow(item.value) + '<br>'
        })
        return result
      }
    },
    legend: {
      data: ['流量使用量'],
      top: 10
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: data.data_axis,
      axisLine: {
        lineStyle: {
          color: '#e6e6e6'
        }
      }
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        formatter: (value) => formatFlow(value)
      },
      axisLine: {
        lineStyle: {
          color: '#e6e6e6'
        }
      },
      splitLine: {
        lineStyle: {
          color: '#f5f5f5'
        }
      }
    },
    series: [
      {
        name: '流量使用量',
        type: 'line',
        stack: 'Total',
        smooth: true,
        lineStyle: {
          width: 3,
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 1, y2: 0,
            colorStops: [
              { offset: 0, color: '#409EFF' },
              { offset: 1, color: '#67C23A' }
            ]
          }
        },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(64, 158, 255, 0.3)' },
              { offset: 1, color: 'rgba(64, 158, 255, 0.1)' }
            ]
          }
        },
        data: data.data
      }
    ]
  })
}

// 监听图表数据变化
watch(() => chartData, (newData) => {
  if (chart.value) {
    setOptions(newData)
  }
}, {
  deep: true
})

onMounted(async () => {
  getTableData()
  setChartData({...searchInfo.value})
  
  await nextTick()
  initChart()
  
  // 监听窗口大小变化
  const handleResize = () => {
    chart.value?.resize()
  }
  window.addEventListener('resize', handleResize)
  
  onUnmounted(() => {
    window.removeEventListener('resize', handleResize)
    if (chart.value) {
      chart.value.dispose()
      chart.value = null
    }
  })
})

const users = ref([])
</script>
<style scoped>
.page {
  padding: 20px;
  background: #f5f7fa;
  min-height: 100vh;
}

/* 图表区域 */
.charts-section {
  margin-bottom: 24px;
}

.chart-card {
  border-radius: 16px;
  overflow: hidden;
  border: none;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
}

.chart-card:hover {
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
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
  color: #303133;
}

.card-title .el-icon {
  margin-right: 8px;
  font-size: 18px;
  color: #409eff;
}

.chart-container {
  height: 360px;
  width: 100%;
}

.trend-chart {
  background: linear-gradient(180deg, #f8f9ff 0%, #ffffff 100%);
}

/* 表格区域 */
.gva-table-box {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.table-header h3 {
  margin: 0;
  color: #303133;
  font-size: 18px;
  font-weight: 600;
}

/* 表格单元格样式 */
.user-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.user-avatar {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  font-weight: bold;
  font-size: 12px;
}

.traffic-down {
  color: #67c23a;
  font-weight: 500;
}

.traffic-up {
  color: #e6a23c;
  font-weight: 500;
}

.date-cell {
  color: #909399;
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
    border-radius: 12px;
  }

  .table-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .table-header h3 {
    font-size: 16px;
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
      background: linear-gradient(to left, rgba(255,255,255,0.9), transparent);
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
    background: #1a1a1a;
  }

  .chart-card,
  .gva-table-box {
    background: #2d2d2d;
    color: #e4e7ed;
  }

  .card-title {
    color: #e4e7ed;
  }

  .table-header h3 {
    color: #e4e7ed;
  }

  .table-header {
    border-bottom-color: #4c4d4f;
  }
}

/* 动画效果 */
.chart-card,
.gva-table-box {
  animation: fadeInUp 0.6s ease-out;
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
  background: #fafbfc;
}

/* 分页样式优化 */
.gva-pagination {
  margin-top: 24px;
  display: flex;
  justify-content: center;
}

:deep(.el-pagination) {
  --el-pagination-button-color: #606266;
  --el-pagination-hover-color: #409eff;
}
</style>