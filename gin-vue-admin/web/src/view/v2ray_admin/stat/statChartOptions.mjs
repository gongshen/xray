import { formatFlow } from './statTraffic.mjs'

const axisLine = {
  lineStyle: {
    color: '#e6e6e6',
  },
}

const splitLine = {
  lineStyle: {
    color: '#f5f5f5',
  },
}

function formatTooltip(params = []) {
  if (!params.length) {
    return ''
  }

  let result = `${params[0].name}<br>`
  params.forEach((item) => {
    result += `${item.marker} ${item.seriesName} : ${formatFlow(item.value)}<br>`
  })
  return result
}

export function buildTrendChartOptions(data = {}) {
  return {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
        label: {
          backgroundColor: '#6a7985',
        },
      },
      formatter: formatTooltip,
    },
    legend: {
      data: ['流量使用量'],
      top: 10,
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: data.data_axis || [],
      axisLine,
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        formatter: (value) => formatFlow(value),
      },
      axisLine,
      splitLine,
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
            x: 0,
            y: 0,
            x2: 1,
            y2: 0,
            colorStops: [
              { offset: 0, color: '#409EFF' },
              { offset: 1, color: '#67C23A' },
            ],
          },
        },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(64, 158, 255, 0.3)' },
              { offset: 1, color: 'rgba(64, 158, 255, 0.1)' },
            ],
          },
        },
        data: data.data || [],
      },
    ],
  }
}

export function buildRankChartOptions(data = {}) {
  return {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow',
      },
      formatter: formatTooltip,
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true,
    },
    xAxis: {
      type: 'value',
      axisLabel: {
        formatter: (value) => formatFlow(value),
      },
      axisLine,
      splitLine,
    },
    yAxis: {
      type: 'category',
      data: data.rank_axis || [],
      axisLine,
    },
    series: [
      {
        name: '流量使用量',
        type: 'bar',
        itemStyle: {
          borderRadius: [0, 6, 6, 0],
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 1,
            y2: 0,
            colorStops: [
              { offset: 0, color: '#FF6B6B' },
              { offset: 0.5, color: '#4ECDC4' },
              { offset: 1, color: '#45B7D1' },
            ],
          },
        },
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)',
          },
        },
        data: data.rank || [],
      },
    ],
  }
}
