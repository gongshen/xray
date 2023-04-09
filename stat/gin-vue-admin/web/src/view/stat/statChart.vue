<template>
  <div class="dashboard-line-box">
    <div ref="echart" class="dashboard-line"></div>
  </div>
</template>
<script setup>
import * as echarts from 'echarts'
import {nextTick, onMounted, onUnmounted, ref, shallowRef, watch} from 'vue'
import { useChartData } from './common'

const chart = shallowRef(null)
const echart = ref(null)

const initChart = () => {
  chart.value = echarts.init(echart.value)
  setOptions()
}

const setOptions = () => {
  const chartData = useChartData()
  chart.value.setOption({
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow',
      },
    },
    grid: {
      left: '40',
      right: '20',
      top: '20',
      bottom: '20',
      containLabel: true,
    },
    xAxis: {
      data: chartData.data_axis,
      axisTick: {
        show: false,
      },
      axisLine: {
        show: false,
      },
      z: 10,
    },
    yAxis: {
      axisLine: {
        show: false,
      },
      axisTick: {
        show: false,
      },
      axisLabel: {
        textStyle: {
          color: '#999',
        },
      },
    },
    dataZoom: [
      {
        type: 'inside',
      },
    ],
    series: [
      {
        type: 'bar',
        barWidth: '40%',
        itemStyle: {
          borderRadius: [5, 5, 0, 0],
          color: '#188df0',
        },
        emphasis: {
          itemStyle: {
            color: '#188df0',
          },
        },
        data: chartData.data,
      },
    ],
  })
}

onMounted(async() => {
  await nextTick()
  initChart()
})

onUnmounted(() => {
  if (!chart.value) {
    return
  }
  chart.value.dispose()
  chart.value = null
})

watch(useChartData, () => {
  console.log("3333")
  setOptions()
})

</script>
<style lang="scss" scoped>
.dashboard-line-box {
.dashboard-line {
  background-color: #fff;
  height: 360px;
  width: 100%;
}
}
</style>
