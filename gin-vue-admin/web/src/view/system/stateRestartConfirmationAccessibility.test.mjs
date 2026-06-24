import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./state.vue', import.meta.url), 'utf8')

const restartButton = source.match(/<el-button[^>]*@click="restartVPS"[^>]*>/)
const restartConfirmation = source.match(/ElMessageBox\.confirm\([\s\S]*?restartVPSApi/)

assert.ok(restartButton, 'server restart button should be rendered')
assert.match(restartButton[0], /type="danger"/, 'server restart button should use danger styling')
assert.match(restartButton[0], /aria-label="重启当前服务器"/, 'server restart button should have a specific accessible name')

assert.ok(restartConfirmation, 'server restart should ask for confirmation')
assert.match(restartConfirmation[0], /'确定要重启当前服务器吗？重启可能需要几分钟时间。'/, 'server restart confirmation should name the target and duration risk')
assert.match(restartConfirmation[0], /'重启服务器'/, 'server restart confirmation should use a specific title')
assert.match(restartConfirmation[0], /confirmButtonText:\s*'重启'/, 'server restart confirmation should use action-specific text')
assert.match(restartConfirmation[0], /cancelButtonText:\s*'取消'/, 'server restart confirmation should keep a cancel action')
assert.match(restartConfirmation[0], /type:\s*'warning'/, 'server restart confirmation should keep warning intent')

console.log('stateRestartConfirmationAccessibility tests passed')
