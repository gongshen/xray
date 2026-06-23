import { formatFlow } from './statTraffic.mjs'
import { chartPalette } from '../../../style/designTokens.mjs'

const axisLine = {
  lineStyle: {
    color: chartPalette.axisLine,
  },
}

const splitLine = {
  lineStyle: {
    color: chartPalette.splitLine,
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

function normalizeAxisValue(value) {
  const text = String(value)
  if (/^\d{8}$/.test(text)) {
    return `${text.slice(0, 4)}-${text.slice(4, 6)}-${text.slice(6, 8)}`
  }
  return value
}

function formatAxisLabel(value) {
  const text = String(value)
  if (/^\d{4}-\d{2}-\d{2}$/.test(text)) {
    return text.slice(5)
  }
  return value
}

export function buildTrendChartOptions(data = {}) {
  return {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
        label: {
          backgroundColor: chartPalette.axisPointerLabel,
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
      data: (data.data_axis || []).map(normalizeAxisValue),
      axisLine,
      axisLabel: {
        rotate: 45,
        formatter: formatAxisLabel,
      },
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
              { offset: 0, color: chartPalette.trendLineStart },
              { offset: 1, color: chartPalette.trendLineEnd },
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
              { offset: 0, color: chartPalette.trendAreaStart },
              { offset: 1, color: chartPalette.trendAreaEnd },
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
              { offset: 0, color: chartPalette.rankStart },
              { offset: 0.5, color: chartPalette.rankMid },
              { offset: 1, color: chartPalette.rankEnd },
            ],
          },
        },
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: chartPalette.emphasisShadow,
          },
        },
        data: data.rank || [],
      },
    ],
  }
}
