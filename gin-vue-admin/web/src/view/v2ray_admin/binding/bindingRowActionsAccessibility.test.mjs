import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./binding.vue', import.meta.url), 'utf8')
const rowActions = source.match(/<div[^>]*class="table-row-actions"[^>]*>/)

assert.ok(rowActions, 'binding table should render a row actions group')
assert.match(rowActions[0], /role="group"/, 'binding row actions should expose grouped controls')
assert.match(rowActions[0], /aria-label="绑定行操作"/, 'binding row actions should have an explicit group name')

console.log('bindingRowActionsAccessibility tests passed')
