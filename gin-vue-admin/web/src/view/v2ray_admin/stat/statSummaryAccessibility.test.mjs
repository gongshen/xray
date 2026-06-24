import assert from 'node:assert/strict'
import fs from 'node:fs'

const files = [
  'src/view/v2ray_admin/stat/stat.vue',
  'src/view/v2ray/stat/stat.vue',
]

for (const file of files) {
  const source = fs.readFileSync(file, 'utf8')
  const summary = source.match(/<div[^>]*class="table-summary"[^>]*>/)

  assert.ok(summary, `${file} should render a traffic summary`)
  assert.match(summary[0], /role="status"/, `${file} traffic summary should expose status semantics`)
  assert.match(summary[0], /aria-live="polite"/, `${file} traffic summary should announce updates politely`)
  assert.match(summary[0], /aria-label="流量统计摘要"/, `${file} traffic summary should have an explicit accessible name`)
}

console.log('statSummaryAccessibility tests passed')
