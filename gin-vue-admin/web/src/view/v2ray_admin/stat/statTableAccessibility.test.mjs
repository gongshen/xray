import assert from 'node:assert/strict'
import fs from 'node:fs'

const expectations = [
  {
    file: 'src/view/v2ray_admin/stat/stat.vue',
    busyBinding: ':aria-busy="tableLoading"',
  },
  {
    file: 'src/view/v2ray/stat/stat.vue',
    busyBinding: ':aria-busy="loading"',
  },
]

for (const { file, busyBinding } of expectations) {
  const source = fs.readFileSync(file, 'utf8')
  const tableMatch = source.match(/<el-table[\s\S]*?>/)
  assert.ok(tableMatch, `${file} should render a stat table`)
  assert.match(tableMatch[0], /class="stat-table"/, `${file} stat table should use the shared table class`)
  assert.match(tableMatch[0], new RegExp(busyBinding.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')), `${file} stat table should expose loading state with aria-busy`)
}

console.log('statTableAccessibility tests passed')
