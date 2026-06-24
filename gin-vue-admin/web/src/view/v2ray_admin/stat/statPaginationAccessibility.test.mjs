import assert from 'node:assert/strict'
import fs from 'node:fs'

const files = [
  'src/view/v2ray_admin/stat/stat.vue',
  'src/view/v2ray/stat/stat.vue',
]

for (const file of files) {
  const source = fs.readFileSync(file, 'utf8')
  const paginationRegion = source.match(/<div[^>]*class="gva-pagination"[^>]*>/)

  assert.ok(paginationRegion, `${file} should render a pagination region`)
  assert.match(paginationRegion[0], /role="navigation"/, `${file} pagination should expose navigation semantics`)
  assert.match(paginationRegion[0], /aria-label="流量记录分页"/, `${file} pagination should have an explicit navigation name`)
}

console.log('statPaginationAccessibility tests passed')
