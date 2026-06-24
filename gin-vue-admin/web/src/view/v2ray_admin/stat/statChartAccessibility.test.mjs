import assert from 'node:assert/strict'
import fs from 'node:fs'

const expectations = [
  {
    file: 'src/view/v2ray_admin/stat/stat.vue',
    charts: [
      { className: 'trend-chart', labelBinding: 'trafficTrendChartLabel' },
      { className: 'rank-chart', labelBinding: 'trafficRankChartLabel' },
    ],
  },
  {
    file: 'src/view/v2ray/stat/stat.vue',
    charts: [
      { className: 'trend-chart', labelBinding: 'trafficTrendChartLabel' },
    ],
  },
]

for (const { file, charts } of expectations) {
  const source = fs.readFileSync(file, 'utf8')

  for (const { className, labelBinding } of charts) {
    const chartPattern = new RegExp(
      `<div[^>]*class="[^"]*chart-container[^"]*${className}[^"]*"[^>]*>`,
      's'
    )
    const chartMatch = source.match(chartPattern)
    assert.ok(chartMatch, `${file} should render ${className}`)
    assert.match(chartMatch[0], /role="img"/, `${file} ${className} should expose image semantics`)
    assert.match(chartMatch[0], new RegExp(`:aria-label="${labelBinding}"`), `${file} ${className} should bind a date-aware accessible label`)
    assert.match(
      source,
      new RegExp('const\\s+' + labelBinding + '\\s*=\\s*computed\\(\\(\\)\\s*=>\\s*`[^`]+\\$\\{dateRangeText\\.value\\}`'),
      `${file} ${labelBinding} should include the selected date range`
    )
  }
}

console.log('statChartAccessibility tests passed')
