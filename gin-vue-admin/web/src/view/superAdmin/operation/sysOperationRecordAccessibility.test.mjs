import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./sysOperationRecord.vue', import.meta.url), 'utf8')

const methodInput = source.match(/<el-input[^>]*v-model="searchInfo\.method"[^>]*>/)
const pathInput = source.match(/<el-input[^>]*v-model="searchInfo\.path"[^>]*>/)
const statusInput = source.match(/<el-input[^>]*v-model="searchInfo\.status"[^>]*>/)

assert.ok(methodInput, 'operation log method search input should be rendered')
assert.match(methodInput[0], /aria-label="请求方法筛选"/, 'operation log method search input should have an accessible name')

assert.ok(pathInput, 'operation log path search input should be rendered')
assert.match(pathInput[0], /aria-label="请求路径筛选"/, 'operation log path search input should have an accessible name')

assert.ok(statusInput, 'operation log status search input should be rendered')
assert.match(statusInput[0], /aria-label="结果状态码筛选"/, 'operation log status search input should have an accessible name')

console.log('sysOperationRecordAccessibility tests passed')
