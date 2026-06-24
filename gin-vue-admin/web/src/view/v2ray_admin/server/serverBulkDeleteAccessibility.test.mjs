import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./server.vue', import.meta.url), 'utf8')

const triggerButton = source.match(/<el-button[^>]*@click="deleteVisible = true"[^>]*>/)
const cancelButton = source.match(/<el-button[^>]*@click="deleteVisible = false"[^>]*>/)
const confirmButton = source.match(/<el-button[^>]*@click="onDelete"[^>]*>/)

assert.ok(triggerButton, 'server bulk delete trigger should be rendered')
assert.match(triggerButton[0], /aria-label="批量删除服务器"/, 'server bulk delete trigger should have a specific accessible name')

assert.ok(cancelButton, 'server bulk delete cancel button should be rendered')
assert.match(cancelButton[0], /aria-label="取消批量删除服务器"/, 'server bulk delete cancel button should have a specific accessible name')

assert.ok(confirmButton, 'server bulk delete confirm button should be rendered')
assert.match(confirmButton[0], /aria-label="确认批量删除服务器"/, 'server bulk delete confirm button should have a specific accessible name')
assert.match(confirmButton[0], /type="danger"/, 'server bulk delete confirm button should use danger styling')

console.log('serverBulkDeleteAccessibility tests passed')
