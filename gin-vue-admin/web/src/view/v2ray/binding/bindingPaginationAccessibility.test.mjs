import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./binding.vue', import.meta.url), 'utf8')
const paginationRegion = source.match(/<div[^>]*class="gva-pagination"[^>]*>/)

assert.ok(paginationRegion, 'v2ray binding list should render pagination')
assert.match(paginationRegion[0], /role="navigation"/, 'v2ray binding pagination should expose navigation semantics')
assert.match(paginationRegion[0], /aria-label="绑定列表分页"/, 'v2ray binding pagination should have an explicit navigation name')

console.log('v2ray bindingPaginationAccessibility tests passed')
