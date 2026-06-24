import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./server.vue', import.meta.url), 'utf8')
const deleteConfirmation = source.match(/ElMessageBox\.confirm\([\s\S]*?deleteServerFunc\(row\)/)

assert.ok(deleteConfirmation, 'server row delete should ask for confirmation')
assert.match(deleteConfirmation[0], /'确定要删除该服务器吗\?'/, 'server row delete confirmation should name the target')
assert.match(deleteConfirmation[0], /'删除服务器'/, 'server row delete confirmation should use a specific title')
assert.match(deleteConfirmation[0], /confirmButtonText:\s*'删除'/, 'server row delete confirmation should use destructive action text')
assert.match(deleteConfirmation[0], /cancelButtonText:\s*'取消'/, 'server row delete confirmation should keep a cancel action')
assert.match(deleteConfirmation[0], /type:\s*'warning'/, 'server row delete confirmation should keep warning intent')

console.log('serverDeleteConfirmationAccessibility tests passed')
