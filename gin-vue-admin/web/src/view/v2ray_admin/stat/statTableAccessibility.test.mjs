import assert from 'node:assert/strict'
import fs from 'node:fs'

const expectations = [
  {
    file: 'src/view/v2ray_admin/stat/stat.vue',
    busyBinding: ':aria-busy="tableLoading"',
    label: '详细流量记录表',
  },
  {
    file: 'src/view/v2ray/stat/stat.vue',
    busyBinding: ':aria-busy="loading"',
    label: '详细流量记录表',
  },
]

for (const { file, busyBinding, label } of expectations) {
  const source = fs.readFileSync(file, 'utf8')
  const tableMatch = source.match(/<el-table[\s\S]*?>/)
  assert.ok(tableMatch, `${file} should render a stat table`)
  assert.match(tableMatch[0], /class="stat-table"/, `${file} stat table should use the shared table class`)
  assert.match(tableMatch[0], new RegExp(busyBinding.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')), `${file} stat table should expose loading state with aria-busy`)
  assert.match(tableMatch[0], new RegExp(`aria-label="${label}"`), `${file} stat table should have an explicit accessible name`)
}

console.log('statTableAccessibility tests passed')
