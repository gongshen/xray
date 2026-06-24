import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./binding.vue', import.meta.url), 'utf8')

const formMatch = source.match(/<el-form[\s\S]*class="demo-form-inline"[\s\S]*?>/)
const startPicker = source.match(/<el-date-picker[^>]*v-model="searchInfo\.startCreatedAt"[^>]*>/)
const endPicker = source.match(/<el-date-picker[^>]*v-model="searchInfo\.endCreatedAt"[^>]*>/)
const actionGroup = source.match(/<el-form-item[^>]*role="group"[\s\S]*?>/)

assert.ok(formMatch, 'v2ray binding search form should be rendered')
assert.match(formMatch[0], /role="search"/, 'v2ray binding search form should use a search landmark')
assert.match(formMatch[0], /aria-label="用户绑定筛选"/, 'v2ray binding search form should have an explicit landmark name')

assert.ok(startPicker, 'v2ray binding start date picker should be rendered')
assert.ok(endPicker, 'v2ray binding end date picker should be rendered')
assert.match(startPicker[0], /aria-label="绑定开始日期"/, 'v2ray binding start date picker should have an accessible name')
assert.match(endPicker[0], /aria-label="绑定结束日期"/, 'v2ray binding end date picker should have an accessible name')

assert.ok(actionGroup, 'v2ray binding search action group should be rendered')
assert.match(actionGroup[0], /aria-label="绑定查询操作"/, 'v2ray binding search action group should have an explicit group name')

console.log('v2ray bindingSearchAccessibility tests passed')
