export const uiColors = Object.freeze({
  brandPrimary: '#1e40af',
  brandSecondary: '#3b82f6',
  accent: '#f59e0b',
  pageBg: '#f8fafc',
  panelBg: '#ffffff',
  panelMutedBg: '#f7f8fa',
  chartTrendBg: '#f8f9ff',
  chartRankBg: '#fff8f0',
  textStrong: '#303133',
  textRegular: '#606266',
  textMuted: '#475569',
  borderSubtle: '#ebeef5',
  borderMuted: '#f0f0f0',
  divider: '#dcdfe6',
  trafficDown: '#67c23a',
  trafficUp: '#f59e0b',
  danger: '#f56c6c',
  darkPageBg: '#1a1a1a',
  darkPanelBg: '#2d2d2d',
  darkText: '#e4e7ed',
  darkBorder: '#4c4d4f',
  chartAxis: '#e6e6e6',
  chartSplit: '#f5f5f5',
  chartPointer: '#64748b',
  rankDanger: '#ef4444',
  rankCyan: '#14b8a6',
  rankBlue: '#3b82f6',
})

export const shadows = Object.freeze({
  panel: '0 4px 20px rgba(15, 23, 42, 0.08)',
  panelHover: '0 8px 30px rgba(15, 23, 42, 0.12)',
})

export const progressThresholdColors = Object.freeze([
  { color: uiColors.trafficDown, percentage: 20 },
  { color: uiColors.trafficUp, percentage: 40 },
  { color: uiColors.danger, percentage: 80 },
])

export const chartPalette = Object.freeze({
  axisLine: uiColors.chartAxis,
  splitLine: uiColors.chartSplit,
  axisPointerLabel: uiColors.chartPointer,
  trendLineStart: uiColors.brandSecondary,
  trendLineEnd: uiColors.trafficDown,
  trendAreaStart: 'rgba(59, 130, 246, 0.3)',
  trendAreaEnd: 'rgba(59, 130, 246, 0.1)',
  rankStart: uiColors.rankDanger,
  rankMid: uiColors.rankCyan,
  rankEnd: uiColors.rankBlue,
  emphasisShadow: 'rgba(15, 23, 42, 0.35)',
})
