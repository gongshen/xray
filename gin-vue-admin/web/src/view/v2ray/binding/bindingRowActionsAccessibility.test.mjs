import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./binding.vue', import.meta.url), 'utf8')
const rowActions = source.match(/<div[^>]*class="table-row-actions"[^>]*>/)

assert.ok(rowActions, 'v2ray binding table should render a row actions group')
assert.match(rowActions[0], /role="group"/, 'v2ray binding row actions should expose grouped controls')
assert.match(rowActions[0], /aria-label="绑定行操作"/, 'v2ray binding row actions should have an explicit group name')

console.log('v2ray bindingRowActionsAccessibility tests passed')
