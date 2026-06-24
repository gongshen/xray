import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync('src/view/v2ray_admin/stat/stat.vue', 'utf8')

const userFilter = source.match(/<el-select[^>]*v-model="searchInfo\.tag"[^>]*>/)
const serverFilter = source.match(/<el-select[^>]*v-model="searchInfo\.server_ip"[^>]*>/)

assert.ok(userFilter, 'admin stat page should render a user filter select')
assert.ok(serverFilter, 'admin stat page should render a server filter select')
assert.match(userFilter[0], /aria-label="统计用户筛选"/, 'user filter should have an explicit accessible name')
assert.match(serverFilter[0], /aria-label="统计服务器筛选"/, 'server filter should have an explicit accessible name')

console.log('statFilterAccessibility tests passed')
