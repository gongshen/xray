import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./server.vue', import.meta.url), 'utf8')

const formMatch = source.match(/<el-form[\s\S]*class="demo-form-inline"[\s\S]*?>/)
const startPicker = source.match(/<el-date-picker[^>]*v-model="searchInfo\.startCreatedAt"[^>]*>/)
const endPicker = source.match(/<el-date-picker[^>]*v-model="searchInfo\.endCreatedAt"[^>]*>/)
const ipInput = source.match(/<el-input[^>]*v-model="searchInfo\.ip"[^>]*>/)
const actionGroup = source.match(/<el-form-item[^>]*role="group"[\s\S]*?>/)

assert.ok(formMatch, 'server search form should be rendered')
assert.match(formMatch[0], /role="search"/, 'server search form should use a search landmark')
assert.match(formMatch[0], /aria-label="服务器筛选"/, 'server search form should have an explicit landmark name')

assert.ok(startPicker, 'server start date picker should be rendered')
assert.ok(endPicker, 'server end date picker should be rendered')
assert.match(startPicker[0], /aria-label="服务器创建开始日期"/, 'server start date picker should have an accessible name')
assert.match(endPicker[0], /aria-label="服务器创建结束日期"/, 'server end date picker should have an accessible name')

assert.ok(ipInput, 'server IP filter should be rendered')
assert.match(ipInput[0], /aria-label="服务器 IP 筛选"/, 'server IP filter should have an accessible name')

assert.ok(actionGroup, 'server search action group should be rendered')
assert.match(actionGroup[0], /aria-label="服务器查询操作"/, 'server search action group should have an explicit group name')

console.log('serverSearchAccessibility tests passed')
