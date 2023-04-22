<template>
  <div className="dashboard-line-box">
    <div className="dashboard-line-title">
      总流量: {{ formatFlow(chartData.total) }}
    </div>
    <div ref="echart" className="dashboard-line"></div>
  </div>
  <div className="dashboard-line-box">
    <div className="dashboard-line-title">
      流量排行榜
    </div>
    <div ref="rank_echart" className="dashboard-line"></div>
  </div>
</template>
<script setup>
import * as echarts from 'echarts'
import {nextTick, onMounted, onUnmounted, ref, shallowRef, watch} from 'vue'
import {useChartData} from './common'

const chart = shallowRef(null)
const rank_chart = shallowRef(null)
const echart = ref(null)
const rank_echart = ref(null)
const chartData = useChartData()


const initChart = () => {
  chart.value = echarts.init(echart.value)
  rank_chart.value = echarts.init(rank_echart.value)
  setOptions(chartData)
}

const formatFlow = (value) => {
  if (value >= 1024 * 1024 * 1024) { // 大于等于 1G 显示 G 后缀
    return (value / (1024 * 1024 * 1024)).toFixed(1) + ' G'
  } else if (value >= 1024 * 1024) { // 大于等于 1M 显示 M 后缀
    return (value / (1024 * 1024)).toFixed(1) + ' M'
  } else if (value >= 1024) { // 大于等于 1k 显示 k 后缀
    return (value / 1024).toFixed(1) + ' K'
  } else { // 小于 1k 直接返回原数值
    return value.toFixed(1)
  }
}

const setOptions = (data) => {
  rank_chart.value.setOption({
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow',
      },
      formatter: (params) => {
        let result = params[0].name + '<br>'
        params.forEach(item => {
          result += item.marker + ' ' + item.seriesName + ' : ' + formatFlow(item.value) + '<br>' // 添加 B/s 后缀
        })
        return result
      }
    },
    grid: {
      left: '1%',
      right: '1%',
      top: '1%',
      bottom: '1%',
      containLabel: true,
    },
    xAxis: {
      data: data.rank_axis,
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
        color: '#999',
        show: false,
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
        barWidth: '20%',
        itemStyle: {
          borderRadius: [5, 5, 0, 0],
          color: '#188df0',
        },
        emphasis: {
          itemStyle: {
            color: '#188df0',
          },
        },
        data: data.rank,
      },
    ],
  });
  chart.value.setOption({
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow',
      },
      formatter: (params) => {
        let result = params[0].name + '<br>'
        params.forEach(item => {
          result += item.marker + ' ' + item.seriesName + ' : ' + formatFlow(item.value) + '<br>' // 添加 B/s 后缀
        })
        return result
      }
    },
    grid: {
      left: '1%',
      right: '1%',
      top: '1%',
      bottom: '1%',
      containLabel: true,
    },
    xAxis: {
      data: data.data_axis,
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
        color: '#999',
        show: false,
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
        barWidth: '20%',
        itemStyle: {
          borderRadius: [5, 5, 0, 0],
          color: '#188df0',
        },
        emphasis: {
          itemStyle: {
            color: '#188df0',
          },
        },
        data: data.data,
      },
    ],
  })
}

onMounted(async () => {
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

watch(() => chartData, (newData) => {
  setOptions(newData)
}, {
  deep: true
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