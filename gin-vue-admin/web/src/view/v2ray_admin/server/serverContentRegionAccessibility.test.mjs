import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./server.vue', import.meta.url), 'utf8')
const tableRegion = source.match(/<div[^>]*class="gva-table-box"[^>]*>/)

assert.ok(tableRegion, 'server list should render a table region')
assert.match(tableRegion[0], /role="region"/, 'server list table area should expose region semantics')
assert.match(tableRegion[0], /aria-label="服务器列表明细"/, 'server list table area should have an explicit region name')

console.log('serverContentRegionAccessibility tests passed')
