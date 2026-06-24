import assert from 'node:assert/strict'
import fs from 'node:fs'

const files = [
  'src/view/v2ray_admin/stat/stat.vue',
  'src/view/v2ray/stat/stat.vue',
]

for (const file of files) {
  const source = fs.readFileSync(file, 'utf8')
  const searchRegion = source.match(/<div[^>]*class="gva-search-box"[^>]*>/)

  assert.ok(searchRegion, `${file} should render a search filter region`)
  assert.match(searchRegion[0], /role="search"/, `${file} search filters should expose search landmark semantics`)
  assert.match(searchRegion[0], /aria-label="流量统计筛选"/, `${file} search filters should have an explicit landmark name`)
}

console.log('statSearchRegionAccessibility tests passed')
