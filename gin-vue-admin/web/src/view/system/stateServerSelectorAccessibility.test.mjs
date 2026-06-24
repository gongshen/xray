import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync('src/view/system/state.vue', 'utf8')
const selectMatch = source.match(/<el-select[\s\S]*class="server-select"[\s\S]*?>/)

assert.ok(selectMatch, 'server selector should be rendered')
assert.match(selectMatch[0], /aria-label="代理服务器"/, 'server selector should have an explicit accessible name')

console.log('stateServerSelectorAccessibility tests passed')
