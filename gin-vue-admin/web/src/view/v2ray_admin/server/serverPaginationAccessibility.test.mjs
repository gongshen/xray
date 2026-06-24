import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./server.vue', import.meta.url), 'utf8')
const paginationRegion = source.match(/<div[^>]*class="gva-pagination"[^>]*>/)

assert.ok(paginationRegion, 'server list should render pagination')
assert.match(paginationRegion[0], /role="navigation"/, 'server pagination should expose navigation semantics')
assert.match(paginationRegion[0], /aria-label="服务器列表分页"/, 'server pagination should have an explicit navigation name')

console.log('serverPaginationAccessibility tests passed')
