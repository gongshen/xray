import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./server.vue', import.meta.url), 'utf8')
const toolbar = source.match(/<div[^>]*class="gva-btn-list"[^>]*>/)

assert.ok(toolbar, 'server list should render a toolbar')
assert.match(toolbar[0], /role="group"/, 'server list toolbar should expose grouped controls')
assert.match(toolbar[0], /aria-label="服务器列表操作"/, 'server list toolbar should have an explicit group name')

console.log('serverToolbarAccessibility tests passed')
