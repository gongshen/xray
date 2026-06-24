import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./binding.vue', import.meta.url), 'utf8')
const tableMatch = source.match(/<el-table[\s\S]*?>/)

assert.ok(tableMatch, 'binding list table should be rendered')
assert.match(tableMatch[0], /v-loading="tableLoading"/, 'binding list table should show loading feedback')
assert.match(tableMatch[0], /:aria-busy="tableLoading"/, 'binding list table should expose loading state to assistive tech')
assert.match(tableMatch[0], /aria-label="用户绑定列表"/, 'binding list table should have an explicit accessible name')
assert.match(tableMatch[0], /:empty-text="tableLoading \? '加载中\.\.\.' : '暂无数据'"/, 'binding list table should clarify empty text while loading')

console.log('bindingTableLoadingAccessibility tests passed')
