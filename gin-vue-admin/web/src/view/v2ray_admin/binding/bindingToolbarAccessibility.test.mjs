import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./binding.vue', import.meta.url), 'utf8')
const toolbar = source.match(/<div[^>]*class="gva-btn-list"[^>]*>/)

assert.ok(toolbar, 'binding list should render a toolbar')
assert.match(toolbar[0], /role="group"/, 'binding list toolbar should expose grouped controls')
assert.match(toolbar[0], /aria-label="绑定列表操作"/, 'binding list toolbar should have an explicit group name')

console.log('bindingToolbarAccessibility tests passed')
