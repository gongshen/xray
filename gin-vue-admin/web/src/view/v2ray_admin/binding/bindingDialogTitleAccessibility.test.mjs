import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./binding.vue', import.meta.url), 'utf8')
const formDialog = source.match(/<el-dialog[^>]*v-model="dialogFormVisible"[^>]*>/)

assert.ok(formDialog, 'binding form dialog should be rendered')
assert.match(formDialog[0], /title="绑定信息"/, 'binding form dialog should have a specific title')

console.log('bindingDialogTitleAccessibility tests passed')
