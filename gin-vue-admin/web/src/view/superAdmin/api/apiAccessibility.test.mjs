import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./api.vue', import.meta.url), 'utf8')

const pathInput = source.match(/<el-input[^>]*v-model="searchInfo\.path"[^>]*>/)
const descriptionInput = source.match(/<el-input[^>]*v-model="searchInfo\.description"[^>]*>/)
const apiGroupInput = source.match(/<el-input[^>]*v-model="searchInfo\.apiGroup"[^>]*>/)
const methodSelect = source.match(/<el-select[^>]*v-model="searchInfo\.method"[^>]*>/)

assert.ok(pathInput, 'API path search input should be rendered')
assert.match(pathInput[0], /aria-label="API 路径筛选"/, 'API path search input should have an accessible name')

assert.ok(descriptionInput, 'API description search input should be rendered')
assert.match(descriptionInput[0], /aria-label="API 描述筛选"/, 'API description search input should have an accessible name')

assert.ok(apiGroupInput, 'API group search input should be rendered')
assert.match(apiGroupInput[0], /aria-label="API 分组筛选"/, 'API group search input should have an accessible name')

assert.ok(methodSelect, 'API method search select should be rendered')
assert.match(methodSelect[0], /aria-label="请求方法筛选"/, 'API method search select should have an accessible name')

console.log('apiAccessibility tests passed')
