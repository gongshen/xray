import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./server.vue', import.meta.url), 'utf8')
const actionColumn = source.match(/<el-table-column[^>]*label="操作"[^>]*>/)

assert.ok(actionColumn, 'server table action column should use a concise user-facing label')

console.log('serverActionColumnLabel tests passed')
