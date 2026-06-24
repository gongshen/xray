import assert from 'node:assert/strict'
import fs from 'node:fs'

const files = [
  'src/view/v2ray_admin/stat/stat.vue',
  'src/view/v2ray/stat/stat.vue',
]

for (const file of files) {
  const source = fs.readFileSync(file, 'utf8')
  const charts = source.match(/<div[^>]*class="charts-section"[^>]*>/)
  const table = source.match(/<div[^>]*class="gva-table-box"[^>]*>/)

  assert.ok(charts, `${file} should render a charts section`)
  assert.match(charts[0], /role="region"/, `${file} charts section should expose region semantics`)
  assert.match(charts[0], /aria-label="流量统计图表"/, `${file} charts section should have an explicit region name`)

  assert.ok(table, `${file} should render a detail table section`)
  assert.match(table[0], /role="region"/, `${file} detail table section should expose region semantics`)
  assert.match(table[0], /aria-label="流量记录明细"/, `${file} detail table section should have an explicit region name`)
}

console.log('statContentRegionAccessibility tests passed')
