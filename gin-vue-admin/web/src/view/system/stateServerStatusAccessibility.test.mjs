import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync('src/view/system/state.vue', 'utf8')
const statusTags = source.match(/<el-tag[^>]*:type="isOnline \? 'success' : 'danger'"[^>]*>/g) || []

assert.equal(statusTags.length, 2, 'server online status should be shown in the page card and traffic analysis dialog')

for (const tag of statusTags) {
  assert.match(tag, /role="status"/, 'server online status tag should expose status semantics')
  assert.match(tag, /aria-live="polite"/, 'server online status updates should be announced politely')
}

console.log('stateServerStatusAccessibility tests passed')
