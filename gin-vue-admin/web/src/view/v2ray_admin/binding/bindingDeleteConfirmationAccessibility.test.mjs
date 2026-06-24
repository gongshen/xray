import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./binding.vue', import.meta.url), 'utf8')
const deleteConfirmation = source.match(/ElMessageBox\.confirm\([\s\S]*?deleteBindingFunc\(row\)/)

assert.ok(deleteConfirmation, 'binding row delete should ask for confirmation')
assert.match(deleteConfirmation[0], /'确定要删除该绑定吗\?'/, 'binding row delete confirmation should name the target')
assert.match(deleteConfirmation[0], /'删除绑定'/, 'binding row delete confirmation should use a specific title')
assert.match(deleteConfirmation[0], /confirmButtonText:\s*'删除'/, 'binding row delete confirmation should use destructive action text')
assert.match(deleteConfirmation[0], /cancelButtonText:\s*'取消'/, 'binding row delete confirmation should keep a cancel action')
assert.match(deleteConfirmation[0], /type:\s*'warning'/, 'binding row delete confirmation should keep warning intent')

console.log('bindingDeleteConfirmationAccessibility tests passed')
