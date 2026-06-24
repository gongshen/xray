import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync('src/view/system/state.vue', 'utf8')
const style = source.match(/<style scoped>([\s\S]*?)<\/style>/)

assert.ok(style, 'state.vue should keep a plain scoped CSS block')
assert.doesNotMatch(
  style[1],
  /\.server-selector\s*\{\s*\.el-button\s*\{/,
  'plain CSS should not rely on Sass-style nested selectors'
)
assert.match(
  style[1],
  /\.server-selector\s+:deep\(\.el-button\)\s*\{[\s\S]*width:\s*100%;/,
  'mobile server action buttons should use an explicit scoped selector'
)

console.log('stateStyleCompatibility tests passed')
