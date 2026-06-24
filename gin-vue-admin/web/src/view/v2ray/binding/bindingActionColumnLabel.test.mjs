import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./binding.vue', import.meta.url), 'utf8')
const actionColumn = source.match(/<el-table-column[^>]*label="操作"[^>]*>/)

assert.ok(actionColumn, 'v2ray binding table action column should use a concise user-facing label')

console.log('v2ray bindingActionColumnLabel tests passed')
