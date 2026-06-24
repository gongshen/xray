import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./server.vue', import.meta.url), 'utf8')

const restartButton = source.match(/<el-button[^>]*@click="restartXray\(scope\.row\)"[^>]*>/)
const restartConfirmation = source.match(/ElMessageBox\.confirm\([\s\S]*?restartXrayFunc\(row\)/)

assert.ok(restartButton, 'server row restart button should be rendered')
assert.match(restartButton[0], /type="warning"/, 'server row restart button should use warning styling')
assert.match(restartButton[0], /aria-label="重启服务器代理"/, 'server row restart button should have a specific accessible name')

assert.ok(restartConfirmation, 'server proxy restart should ask for confirmation')
assert.match(restartConfirmation[0], /'确定要重启该服务器代理吗\?'/, 'server proxy restart confirmation should name the target')
assert.match(restartConfirmation[0], /'重启代理'/, 'server proxy restart confirmation should use a specific title')
assert.match(restartConfirmation[0], /confirmButtonText:\s*'重启'/, 'server proxy restart confirmation should use action-specific text')
assert.match(restartConfirmation[0], /cancelButtonText:\s*'取消'/, 'server proxy restart confirmation should keep a cancel action')
assert.match(restartConfirmation[0], /type:\s*'warning'/, 'server proxy restart confirmation should keep warning intent')

console.log('serverRestartConfirmationAccessibility tests passed')
