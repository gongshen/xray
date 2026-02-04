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

    <!-- 统计卡片区域 -->
    <div class="stats-overview">
      <el-row :gutter="20">
        <el-col :xs="24" :sm="12" :md="8" :lg="8">
          <div class="stat-card total-traffic">
            <div class="stat-icon">
              <el-icon><TrendCharts /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-title">总流量</div>
              <div class="stat-value">{{ formatFlow(chartData.total) }}</div>
              <div class="stat-desc">累计使用流量</div>
            </div>
          </div>
        </el-col>
        <el-col :xs="24" :sm="12" :md="8" :lg="8">
          <div class="stat-card active-servers">
            <div class="stat-icon">
              <el-icon><Monitor /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-title">活跃服务器</div>
              <div class="stat-value">{{ activeServers }}</div>
              <div class="stat-desc">有流量产生的服务器</div>
            </div>
          </div>
        </el-col>
        <el-col :xs="24" :sm="12" :md="8" :lg="8">
          <div class="stat-card avg-traffic">
            <div class="stat-icon">
              <el-icon><DataAnalysis /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-title">平均流量</div>
              <div class="stat-value">{{ formatFlow(avgTraffic) }}</div>
              <div class="stat-desc">每用户平均使用</div>
            </div>
          </div>
        </el-col>
      </el-row>
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
                <el-tag type="info" size="small">{{ getDateRangeText() }}</el-tag>
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
        <el-tag>共 {{ total }} 条记录</el-tag>
      </div>
      
      <el-table
          ref="multipleTable"
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
import { ref, reactive, shallowRef, onMounted, nextTick, onUnmounted, computed, watch } from 'vue'
import { TrendCharts, Monitor, DataAnalysis, Trophy } from '@element-plus/icons-vue'
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
const monthAgo = new Date(today.getTime() - 30 * 24 * 60 * 60 * 1000) // 1个月前
const searchInfo = ref({
  startCreatedAt: monthAgo.toISOString(),
  endCreatedAt: today.toISOString()
})

// 图表相关
const chart = shallowRef(null)
const rank_chart = shallowRef(null)
const echart = ref(null)
const rank_echart = ref(null)
const chartData = useChartData()

// 统计数据计算
const activeServers = computed(() => {
  const uniqueServers = new Set(tableData.value.map(item => item.server_ip))
  return uniqueServers.size
})

const avgTraffic = computed(() => {
  const uniqueUsers = new Set(tableData.value.map(item => item.tag))
  const activeUsers = uniqueUsers.size
  if (activeUsers === 0) return 0
  return chartData.total / activeUsers
})

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
const getTrafficTagType = (trafficStr) => {
  if (!trafficStr || typeof trafficStr !== 'string') return 'info'
  
  // 解析流量字符串，如 "726.27MB" -> 726.27 * 1024 * 1024
  const match = trafficStr.match(/^([\d.]+)\s*(B|KB|MB|GB)$/i)
  if (!match) return 'info'
  
  const value = parseFloat(match[1])
  const unit = match[2].toUpperCase()
  
  let bytes = value
  switch (unit) {
    case 'KB':
      bytes = value * 1024
      break
    case 'MB':
      bytes = value * 1024 * 1024
      break
    case 'GB':
      bytes = value * 1024 * 1024 * 1024
      break
  }
  
  if (bytes >= 1024 * 1024 * 1024) return 'danger'  // GB级别
  if (bytes >= 100 * 1024 * 1024) return 'warning'  // 100MB以上
  if (bytes >= 10 * 1024 * 1024) return 'success'   // 10MB以上
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
    console.log('=== V2RAY_ADMIN 获取表格数据 ===')
    console.log('请求参数:', { page: page.value, pageSize: pageSize.value, ...searchInfo.value })
    
    const table = await getStatList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
    console.log('API响应:', table)
    
    if (table.code === 0) {
      // 确保数据结构正确
      const list = table.data?.list || []
      tableData.value = list
      total.value = table.data?.total || 0
      page.value = table.data?.page || 1
      pageSize.value = table.data?.pageSize || 10
      
      console.log('表格数据设置完成:', {
        dataLength: tableData.value.length,
        total: total.value,
        sampleData: tableData.value[0] || null
      })
      
      // 如果有数据但表格不显示，可能是响应式问题
      if (list.length > 0) {
        console.log('数据样本:', list[0])
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
  }
}

// 初始化图表
const initChart = () => {
  chart.value = echarts.init(echart.value)
  rank_chart.value = echarts.init(rank_echart.value)
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

  // 流量排行榜
  rank_chart.value?.setOption({
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      },
      formatter: (params) => {
        let result = params[0].name + '<br>'
        params.forEach(item => {
          result += item.marker + ' ' + item.seriesName + ' : ' + formatFlow(item.value) + '<br>'
        })
        return result
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
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
    yAxis: {
      type: 'category',
      data: data.rank_axis,
      axisLine: {
        lineStyle: {
          color: '#e6e6e6'
        }
      }
    },
    series: [
      {
        name: '流量使用量',
        type: 'bar',
        itemStyle: {
          borderRadius: [0, 6, 6, 0],
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 1, y2: 0,
            colorStops: [
              { offset: 0, color: '#FF6B6B' },
              { offset: 0.5, color: '#4ECDC4' },
              { offset: 1, color: '#45B7D1' }
            ]
          }
        },
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        },
        data: data.rank
      }
    ]
  })
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
  if (chart.value && rank_chart.value) {
    setOptions(newData)
  }
}, {
  deep: true
})

onMounted(async () => {
  init()
  getTableData()
  setChartData({...searchInfo.value})
  
  await nextTick()
  initChart()
  
  // 监听窗口大小变化
  const handleResize = () => {
    chart.value?.resize()
    rank_chart.value?.resize()
  }
  window.addEventListener('resize', handleResize)
  
  onUnmounted(() => {
    window.removeEventListener('resize', handleResize)
    if (chart.value) {
      chart.value.dispose()
      chart.value = null
    }
    if (rank_chart.value) {
      rank_chart.value.dispose()
      rank_chart.value = null
    }
  })
})
</script>

<style scoped>
.page {
  padding: 20px;
  background: #f5f7fa;
  min-height: 100vh;
}

/* 统计卡片区域 */
.stats-overview {
  margin-bottom: 24px;
}

.stat-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 16px;
  padding: 24px;
  color: white;
  display: flex;
  align-items: center;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  margin-bottom: 16px;
  position: relative;
  overflow: hidden;
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0.05) 100%);
  pointer-events: none;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
}

.stat-card.total-traffic {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-card.active-servers {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-card.avg-traffic {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.stat-icon {
  font-size: 48px;
  margin-right: 20px;
  opacity: 0.8;
}

.stat-content {
  flex: 1;
}

.stat-title {
  font-size: 14px;
  opacity: 0.9;
  margin-bottom: 8px;
  font-weight: 500;
}

.stat-value {
  font-size: 32px;
  font-weight: bold;
  margin-bottom: 4px;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.stat-desc {
  font-size: 12px;
  opacity: 0.7;
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

.rank-chart {
  background: linear-gradient(180deg, #fff8f0 0%, #ffffff 100%);
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

  .stats-overview {
    margin-bottom: 16px;
  }

  .stat-card {
    padding: 16px;
    margin-bottom: 12px;
  }

  .stat-icon {
    font-size: 36px;
    margin-right: 12px;
  }

  .stat-value {
    font-size: 24px;
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

  /* 统计卡片在移动端堆叠显示 */
  .stat-card {
    flex-direction: column;
    text-align: center;
  }

  .stat-icon {
    margin-right: 0;
    margin-bottom: 12px;
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
.stat-card,
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