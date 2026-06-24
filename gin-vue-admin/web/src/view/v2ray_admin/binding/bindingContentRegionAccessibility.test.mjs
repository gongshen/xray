import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./binding.vue', import.meta.url), 'utf8')
const tableRegion = source.match(/<div[^>]*class="gva-table-box"[^>]*>/)

assert.ok(tableRegion, 'binding list should render a table region')
assert.match(tableRegion[0], /role="region"/, 'binding list table area should expose region semantics')
assert.match(tableRegion[0], /aria-label="绑定列表明细"/, 'binding list table area should have an explicit region name')

console.log('bindingContentRegionAccessibility tests passed')
