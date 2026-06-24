import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync('src/view/system/state.vue', 'utf8')
const emptyState = source.match(/<el-empty[^>]*v-if="!currentServer && !loading"[^>]*>/)

assert.ok(emptyState, 'server empty state should be rendered when no server is selected')
assert.match(emptyState[0], /role="status"/, 'server empty state should expose status semantics')
assert.match(emptyState[0], /aria-live="polite"/, 'server empty state should announce changes politely')
assert.match(emptyState[0], /aria-label="服务器状态空状态"/, 'server empty state should have an explicit accessible name')

console.log('stateEmptyStateAccessibility tests passed')
