import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./sysOperationRecord.vue', import.meta.url), 'utf8')

const methodInput = source.match(/<el-input[^>]*v-model="searchInfo\.method"[^>]*>/)
const pathInput = source.match(/<el-input[^>]*v-model="searchInfo\.path"[^>]*>/)
const statusInput = source.match(/<el-input[^>]*v-model="searchInfo\.status"[^>]*>/)
const requestBodyButton = source.match(/<button[^>]*class="payload-popover-trigger"[^>]*aria-label="查看请求内容"[^>]*>/)
const responseBodyButton = source.match(/<button[^>]*class="payload-popover-trigger"[^>]*aria-label="查看响应内容"[^>]*>/)

assert.ok(methodInput, 'operation log method search input should be rendered')
assert.match(methodInput[0], /aria-label="请求方法筛选"/, 'operation log method search input should have an accessible name')

assert.ok(pathInput, 'operation log path search input should be rendered')
assert.match(pathInput[0], /aria-label="请求路径筛选"/, 'operation log path search input should have an accessible name')

assert.ok(statusInput, 'operation log status search input should be rendered')
assert.match(statusInput[0], /aria-label="结果状态码筛选"/, 'operation log status search input should have an accessible name')

assert.ok(requestBodyButton, 'request payload popover trigger should use a native button')
assert.match(requestBodyButton[0], /type="button"/, 'request payload popover trigger should not submit forms')

assert.ok(responseBodyButton, 'response payload popover trigger should use a native button')
assert.match(responseBodyButton[0], /type="button"/, 'response payload popover trigger should not submit forms')

assert.doesNotMatch(source, /<el-icon[^>]*style="cursor: pointer;"[^>]*><warning \/><\/el-icon>/, 'payload popover triggers should not be cursor-only icons')

console.log('sysOperationRecordAccessibility tests passed')
