import assert from 'node:assert/strict'
import fs from 'node:fs'

const expectations = [
  {
    file: 'src/view/v2ray_admin/stat/stat.vue',
    charts: [
      { className: 'trend-chart', label: '流量趋势图表' },
      { className: 'rank-chart', label: '流量排行榜图表' },
    ],
  },
  {
    file: 'src/view/v2ray/stat/stat.vue',
    charts: [
      { className: 'trend-chart', label: '流量趋势图表' },
    ],
  },
]

for (const { file, charts } of expectations) {
  const source = fs.readFileSync(file, 'utf8')

  for (const { className, label } of charts) {
    const chartPattern = new RegExp(
      `<div[^>]*class="[^"]*chart-container[^"]*${className}[^"]*"[^>]*>`,
      's'
    )
    const chartMatch = source.match(chartPattern)
    assert.ok(chartMatch, `${file} should render ${className}`)
    assert.match(chartMatch[0], /role="img"/, `${file} ${className} should expose image semantics`)
    assert.match(chartMatch[0], new RegExp(`aria-label="${label}"`), `${file} ${className} should have an accessible label`)
  }
}

console.log('statChartAccessibility tests passed')
