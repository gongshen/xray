<template>
  <div class="dashboard-line-box">
    <div class="dashboard-line-title">
      总流量 {{ formatFlow(chartData.total) }}
    </div>
    <div ref="echart" class="dashboard-line" role="img" aria-label="流量趋势图"></div>
  </div>
</template>

<script setup>
import { nextTick, onMounted, onUnmounted, ref, shallowRef, watch } from 'vue'
import { bindWindowEvent } from '@/utils/eventLifecycle.mjs'
import { loadEcharts } from '@/utils/loadEcharts'
import { useChartData } from './common'
import { formatFlow } from '../../v2ray_admin/stat/statTraffic.mjs'
import { buildTrendChartOptions } from '../../v2ray_admin/stat/statChartOptions.mjs'

const chart = shallowRef(null)
const echart = ref(null)
const chartData = useChartData()
let disposeResize = null
let isUnmounted = false

const setOptions = (data) => {
  chart.value?.setOption(buildTrendChartOptions(data))
}

const handleResize = () => {
  chart.value?.resize()
}

const disposeChart = () => {
  if (chart.value) {
    chart.value.dispose()
    chart.value = null
  }
}

const initChart = async() => {
  if (!echart.value) {
    return
  }

  const echarts = await loadEcharts()
  if (!echart.value || isUnmounted) {
    return
  }

  disposeChart()
  chart.value = echarts.init(echart.value)
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
  disposeChart()
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
