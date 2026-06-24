import assert from 'node:assert/strict'
import fs from 'node:fs'

const expectations = [
  {
    file: 'src/view/v2ray_admin/stat/stat.vue',
    loadingRef: 'tableLoading',
  },
  {
    file: 'src/view/v2ray/stat/stat.vue',
    loadingRef: 'loading',
  },
]

for (const { file, loadingRef } of expectations) {
  const source = fs.readFileSync(file, 'utf8')
  const table = source.match(/<el-table[\s\S]*?>/)

  assert.ok(table, `${file} should render a stat table`)
  assert.match(
    table[0],
    new RegExp(`:empty-text="${loadingRef} \\? '加载中\\.\\.\\.' : \\(tableData\\.length === 0 \\? '暂无数据' : ''\\)"`),
    `${file} table should not announce an empty state while loading`
  )
}

console.log('statTableEmptyState tests passed')
