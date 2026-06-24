import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync('src/view/system/state.vue', 'utf8')
const controlsRegion = source.match(/<div[^>]*class="server-selector"[^>]*>/)

assert.ok(controlsRegion, 'server controls should be rendered')
assert.match(controlsRegion[0], /role="group"/, 'server controls should expose group semantics')
assert.match(controlsRegion[0], /aria-label="服务器状态操作"/, 'server controls group should have an explicit accessible name')

console.log('stateServerControlsAccessibility tests passed')
