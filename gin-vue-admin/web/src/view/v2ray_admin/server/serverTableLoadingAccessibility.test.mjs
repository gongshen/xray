import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./server.vue', import.meta.url), 'utf8')
const tableMatch = source.match(/<el-table[\s\S]*?>/)

assert.ok(tableMatch, 'server list table should be rendered')
assert.match(tableMatch[0], /v-loading="tableLoading"/, 'server list table should show loading feedback')
assert.match(tableMatch[0], /:aria-busy="tableLoading"/, 'server list table should expose loading state to assistive tech')
assert.match(tableMatch[0], /aria-label="服务器列表"/, 'server list table should have an explicit accessible name')
assert.match(tableMatch[0], /:empty-text="tableLoading \? '加载中\.\.\.' : '暂无数据'"/, 'server list table should clarify empty text while loading')
assert.match(source, /const tableLoading = ref\(false\)/, 'server list should track table loading state')
assert.match(source, /tableLoading\.value = true/, 'server list should set loading before fetching data')
assert.match(source, /finally\s*\{[\s\S]*tableLoading\.value = false[\s\S]*\}/, 'server list should clear loading after fetching data')

console.log('serverTableLoadingAccessibility tests passed')
