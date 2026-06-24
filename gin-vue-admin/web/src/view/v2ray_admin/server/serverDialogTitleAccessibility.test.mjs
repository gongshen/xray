import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./server.vue', import.meta.url), 'utf8')
const formDialog = source.match(/<el-dialog[^>]*v-model="dialogFormVisible"[^>]*>/)

assert.ok(formDialog, 'server form dialog should be rendered')
assert.match(formDialog[0], /title="服务器信息"/, 'server form dialog should have a specific title')

console.log('serverDialogTitleAccessibility tests passed')
