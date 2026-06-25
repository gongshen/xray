<template>
  <div class="dashboard-line-box">
    <div class="dashboard-line-title">
      总流量 {{ formatFlow(chartData.total) }}
    </div>
    <div ref="echart" class="dashboard-line" role="img" aria-label="流量趋势图"></div>
  </div>
  <div class="dashboard-line-box">
    <div class="dashboard-line-title">
      流量排行榜
    </div>
    <div ref="rankEchart" class="dashboard-line" role="img" aria-label="流量排行榜图"></div>
  </div>
</template>

<script setup>
import { nextTick, onMounted, onUnmounted, ref, shallowRef, watch } from 'vue'
import { bindWindowEvent } from '@/utils/eventLifecycle.mjs'
import { loadEcharts } from '@/utils/loadEcharts'
import { useChartData } from './common'
import { formatFlow } from './statTraffic.mjs'
import {
  buildRankChartOptions,
  buildTrendChartOptions,
} from './statChartOptions.mjs'

const chart = shallowRef(null)
const rankChart = shallowRef(null)
const echart = ref(null)
const rankEchart = ref(null)
const chartData = useChartData()
let disposeResize = null
let isUnmounted = false

const setOptions = (data) => {
  chart.value?.setOption(buildTrendChartOptions(data))
  rankChart.value?.setOption(buildRankChartOptions(data))
}

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

const initChart = async() => {
  if (!echart.value || !rankEchart.value) {
    return
  }

  const echarts = await loadEcharts()
  if (!echart.value || !rankEchart.value || isUnmounted) {
    return
  }

  disposeCharts()
  chart.value = echarts.init(echart.value)
  rankChart.value = echarts.init(rankEchart.value)
  setOptions(chartData)
  disposeResize?.()
  disposeResize = bindWindowEvent(window, 'resize', handleResize)
}

watch(() => chartData, (newData) => {
  setOptions(newData)
}, {
  deep: true,
})

onMounted(async() => {
  await nextTick()
  await initChart()
})

onUnmounted(() => {
  isUnmounted = true
  disposeResize?.()
  disposeResize = null
  disposeCharts()
})
</script>

<style lang="scss" scoped>
.dashboard-line-box {
  .dashboard-line {
    background-color: #fff;
    height: 360px;
    width: 100%;
  }

  .dashboard-line-title {
    font-weight: 600;
    margin-bottom: 12px;
  }
}
</style>
