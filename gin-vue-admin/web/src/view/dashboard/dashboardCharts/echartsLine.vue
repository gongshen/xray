<template>
  <div class="dashboard-line-box">
    <div class="dashboard-line-title">
      访问趋势
    </div>
    <div ref="echart" class="dashboard-line" role="img" aria-label="访问趋势柱状图"></div>
  </div>
</template>
<script setup>
import { nextTick, onMounted, onUnmounted, ref, shallowRef } from 'vue'
import { bindWindowEvent } from '@/utils/eventLifecycle.mjs'
import { loadEcharts } from '@/utils/loadEcharts'
// import 'echarts/theme/macarons'

const dataAxis = [];
for (let i = 1; i < 13; i++) {
  dataAxis.push(`${i}月`)
}

const data = [
  220,
  182,
  191,
  234,
  290,
  330,
  310,
  123,
  442,
  321,
  90,
  149,
];

const chart = shallowRef(null)
const echart = ref(null)
let disposeResize = null
let isUnmounted = false

const initChart = async() => {
  if (!echart.value) {
    return
  }

  const echarts = await loadEcharts()
  if (!echart.value || isUnmounted) {
    return
  }

  chart.value = echarts.init(echart.value)
  setOptions()
  disposeResize = bindWindowEvent(window, 'resize', () => {
    chart.value?.resize()
  })
}
const setOptions = () => {
  chart.value?.setOption({
    grid: {
      left: '40',
      right: '20',
      top: '20',
      bottom: '20',
    },
    xAxis: {
      data: dataAxis,
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
        data: data,
      },
    ],
  })
}

onMounted(async() => {
  await nextTick()
  await initChart()
})

onUnmounted(() => {
  isUnmounted = true
  disposeResize?.()
  disposeResize = null
  if (!chart.value) {
    return
  }
  chart.value.dispose()
  chart.value = null
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
