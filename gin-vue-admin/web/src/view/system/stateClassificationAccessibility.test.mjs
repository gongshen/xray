import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync('src/view/system/state.vue', 'utf8')
const resultMatch = source.match(/<div\s+v-if="targetClassificationResult"[^>]*class="classification-result"[^>]*>/)

assert.ok(resultMatch, 'target classification result region should be rendered when data exists')
assert.match(resultMatch[0], /role="status"/, 'classification result should expose status semantics')
assert.match(resultMatch[0], /aria-live="polite"/, 'classification result should announce async updates politely')

console.log('stateClassificationAccessibility tests passed')
