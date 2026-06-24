<template>
  <div class="page">
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
        <el-form-item label="创建时间">
          <el-date-picker v-model="searchInfo.startCreatedAt" type="date" placeholder="开始时间" aria-label="统计开始日期" :default-value="monthAgo"></el-date-picker>
          <el-date-picker v-model="searchInfo.endCreatedAt" type="date" placeholder="结束时间" aria-label="统计结束日期" :default-value="today"></el-date-picker>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="onSubmit">查询</el-button>
          <el-button :icon="Refresh" @click="onReset">重置</el-button>
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
                <el-tag type="info" size="small">{{ dateRangeText }}</el-tag>
              </div>
            </template>
            <div ref="echart" class="chart-container trend-chart" role="img" aria-label="流量趋势图表"></div>
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
          class="stat-table"
          tooltip-effect="dark"
          :data="tableData"
          row-key="ID"
          stripe
          v-loading="loading"
          :aria-busy="loading"
          aria-label="详细流量记录表"
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
import { ElMessage } from 'element-plus'
import { computed, ref, shallowRef, onMounted, nextTick, onUnmounted, watch } from 'vue'
import { Refresh, Search, TrendCharts } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { useChartData, setChartData } from "./common"
import { bindWindowEvent } from '@/utils/eventLifecycle.mjs'
import {
  createDefaultTrafficSearchRange,
  formatFlow,
  getDateRangeText as formatDateRangeText,
  getTrafficTagType,
  normalizeTrafficSearchRange,
} from '../../v2ray_admin/stat/statTraffic.mjs'
import {
  buildTrendChartOptions,
} from '../../v2ray_admin/stat/statChartOptions.mjs'
import { normalizeStatTableResponse } from '../../v2ray_admin/stat/statTableData.mjs'

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const loading = ref(false)
const initialSearchRange = createDefaultTrafficSearchRange()
const today = new Date(initialSearchRange.endCreatedAt)
const monthAgo = new Date(initialSearchRange.startCreatedAt)
const searchInfo = ref({ ...initialSearchRange })

// 图表相关
const chart = shallowRef(null)
let disposeResize = null
const echart = ref(null)
const chartData = useChartData()
const dateRangeText = computed(() => formatDateRangeText(searchInfo.value))

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
  loading.value = true
  try {
    const table = await getStatList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
    const tableState = normalizeStatTableResponse(table, { page: page.value, pageSize: pageSize.value })

    if (tableState.ok) {
      tableData.value = tableState.list
      total.value = tableState.total
      page.value = tableState.page
      pageSize.value = tableState.pageSize
    } else {
      console.error('getStatList error:', table.msg)
      ElMessage.error(tableState.message || '获取数据失败')
      tableData.value = tableState.list
      total.value = tableState.total
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
    return
  }

  if (chart.value) {
    chart.value.dispose()
  }

  chart.value = echarts.init(echart.value)
  setOptions(chartData)
}

// 设置图表选项
const setOptions = (data) => {
  chart.value?.setOption(buildTrendChartOptions(data))
}

// 监听图表数据变化
watch(() => chartData, (newData) => {
  if (chart.value) {
    setOptions(newData)
  }
}, {
  deep: true
})

const handleResize = () => {
  chart.value?.resize()
}

const disposeChart = () => {
  if (chart.value) {
    chart.value.dispose()
    chart.value = null
  }
}

onMounted(async () => {
  await getTableData()
  await setChartData({...searchInfo.value})
  await nextTick()

  initChart()
  disposeResize = bindWindowEvent(window, 'resize', handleResize)
})

onUnmounted(() => {
  disposeResize?.()
  disposeResize = null
  disposeChart()
})
</script>
<style scoped lang="scss" src="../../v2ray_admin/stat/statPage.scss"></style>
