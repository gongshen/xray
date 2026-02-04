<template>
  <div class="page">
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
        <el-form-item label="创建时间">
          <el-date-picker v-model="searchInfo.startCreatedAt" type="date" placeholder="开始时间" :default-value="monthAgo"></el-date-picker>
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
          style="width: 100%; min-height: 200px;"
          tooltip-effect="dark"
          :data="tableData"
          row-key="ID"
          stripe
          v-loading="loading"
          :empty-text="loading ? '加载中...' : (tableData.length === 0 ? '暂无数据' : '')"
      >
        <el-table-column type="selection" width="55" />
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
const loading = ref(false)
const today = new Date()
const monthAgo = new Date(today.getTime() - 30 * 24 * 60 * 60 * 1000) // 1个月前
const searchInfo = ref({
  startCreatedAt: monthAgo.toISOString(),
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
  setChartData({...searchInfo.value}).then(() => {
    // 数据加载完成后刷新图表
    setTimeout(() => {
      refreshChart()
    }, 100)
  })
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
  setChartData({...searchInfo.value}).then(() => {
    // 数据加载完成后刷新图表
    setTimeout(() => {
      refreshChart()
    }, 100)
  })
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
  loading.value = true
  try {
    console.log('=== 获取表格数据 ===')
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
  } finally {
    loading.value = false
  }
}

// 初始化图表
const initChart = () => {
  if (!echart.value) {
    console.error('图表容器不存在，无法初始化图表')
    return
  }
  
  // 如果图表已存在，先销毁
  if (chart.value) {
    chart.value.dispose()
  }
  
  chart.value = echarts.init(echart.value)
  console.log('图表实例创建完成')
  
  // 如果已有数据，立即渲染
  if (chartData.data && chartData.data.length > 0) {
    console.log('立即渲染已有数据')
    setOptions(chartData)
  }
}

// 强制刷新图表
const refreshChart = () => {
  console.log('强制刷新图表')
  if (chart.value && chartData.data && chartData.data.length > 0) {
    setOptions(chartData)
  } else {
    console.log('重新初始化图表')
    nextTick(() => {
      initChart()
    })
  }
}

// 设置图表选项
const setOptions = (data) => {
  // 添加调试信息
  console.log('图表数据:', {
    data_axis: data.data_axis,
    data: data.data,
    data_axis_length: data.data_axis?.length,
    data_length: data.data?.length
  })

  // 检查数据是否存在
  if (!data.data_axis || !data.data || data.data_axis.length === 0 || data.data.length === 0) {
    console.warn('图表数据为空或格式不正确')
    return
  }

  // 格式化日期轴数据
  const formattedAxisData = data.data_axis.map(dateNum => {
    const dateStr = dateNum.toString()
    if (dateStr.length === 8) {
      // 格式：YYYYMMDD -> YYYY-MM-DD
      const year = dateStr.substring(0, 4)
      const month = dateStr.substring(4, 6)
      const day = dateStr.substring(6, 8)
      return `${year}-${month}-${day}`
    }
    return dateStr
  })

  console.log('格式化后的日期轴:', formattedAxisData)

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
      data: formattedAxisData,
      axisLine: {
        lineStyle: {
          color: '#e6e6e6'
        }
      },
      axisLabel: {
        rotate: 45, // 旋转标签避免重叠
        formatter: (value) => {
          // 简化日期显示，只显示月-日
          if (value.includes('-')) {
            const parts = value.split('-')
            return `${parts[1]}-${parts[2]}`
          }
          return value
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
  console.log('图表数据变化:', newData)
  if (chart.value && newData) {
    setOptions(newData)
  } else if (!chart.value) {
    console.warn('图表实例不存在，尝试重新初始化')
    nextTick(() => {
      if (echart.value) {
        initChart()
      }
    })
  }
}, {
  deep: true
})

onMounted(async () => {
  console.log('页面挂载，开始初始化')
  
  // 先加载表格数据
  await getTableData()
  
  // 再加载图表数据
  await setChartData({...searchInfo.value})
  
  await nextTick()
  console.log('DOM更新完成，初始化图表')
  
  // 确保图表容器存在
  if (echart.value) {
    initChart()
    console.log('图表初始化完成')
    
    // 延迟刷新图表，确保数据已加载
    setTimeout(() => {
      console.log('延迟刷新图表')
      refreshChart()
    }, 500)
  } else {
    console.error('图表容器不存在')
  }
  
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