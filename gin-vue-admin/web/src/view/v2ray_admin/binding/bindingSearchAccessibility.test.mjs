import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./binding.vue', import.meta.url), 'utf8')

const formMatch = source.match(/<el-form[\s\S]*class="demo-form-inline"[\s\S]*?>/)
const startPicker = source.match(/<el-date-picker[^>]*v-model="searchInfo\.startCreatedAt"[^>]*>/)
const endPicker = source.match(/<el-date-picker[^>]*v-model="searchInfo\.endCreatedAt"[^>]*>/)
const serverSelect = source.match(/<el-select[^>]*v-model="searchInfo\.server_id"[^>]*>/)
const userSelect = source.match(/<el-select[^>]*v-model="searchInfo\.user_id"[^>]*>/)
const actionGroup = source.match(/<el-form-item[^>]*role="group"[\s\S]*?>/)

assert.ok(formMatch, 'binding search form should be rendered')
assert.match(formMatch[0], /role="search"/, 'binding search form should use a search landmark')
assert.match(formMatch[0], /aria-label="绑定筛选"/, 'binding search form should have an explicit landmark name')

assert.ok(startPicker, 'binding start date picker should be rendered')
assert.ok(endPicker, 'binding end date picker should be rendered')
assert.match(startPicker[0], /aria-label="绑定创建开始日期"/, 'binding start date picker should have an accessible name')
assert.match(endPicker[0], /aria-label="绑定创建结束日期"/, 'binding end date picker should have an accessible name')

assert.ok(serverSelect, 'binding server filter should be rendered')
assert.ok(userSelect, 'binding user filter should be rendered')
assert.match(serverSelect[0], /aria-label="绑定服务器筛选"/, 'binding server filter should have an accessible name')
assert.match(userSelect[0], /aria-label="绑定用户筛选"/, 'binding user filter should have an accessible name')

assert.ok(actionGroup, 'binding search action group should be rendered')
assert.match(actionGroup[0], /aria-label="绑定查询操作"/, 'binding search action group should have an explicit group name')

console.log('bindingSearchAccessibility tests passed')
