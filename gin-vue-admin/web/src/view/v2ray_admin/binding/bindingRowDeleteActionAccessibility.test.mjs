import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./binding.vue', import.meta.url), 'utf8')
const deleteButton = source.match(/<el-button[^>]*@click="deleteRow\(scope\.row\)"[^>]*>/)

assert.ok(deleteButton, 'binding row delete button should be rendered')
assert.match(deleteButton[0], /type="danger"/, 'binding row delete button should use danger styling')
assert.match(deleteButton[0], /aria-label="删除绑定"/, 'binding row delete button should have a specific accessible name')

console.log('bindingRowDeleteActionAccessibility tests passed')
