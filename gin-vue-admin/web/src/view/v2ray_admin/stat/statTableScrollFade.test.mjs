import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync('src/view/v2ray_admin/stat/statPage.scss', 'utf8')

assert.doesNotMatch(
  source,
  /rgba\(\s*255\s*,\s*255\s*,\s*255\s*,\s*0\.9\s*\)/,
  'mobile table scroll fade should not hard-code a light-only white background'
)
assert.match(
  source,
  /&::before\s*\{[\s\S]*background:\s*linear-gradient\(to left,\s*var\(--gva-color-panel-bg\),\s*transparent\);/,
  'mobile table scroll fade should use the shared light panel token'
)
assert.match(
  source,
  /@media \(prefers-color-scheme:\s*dark\)[\s\S]*\.gva-table-box::before\s*\{[\s\S]*background:\s*linear-gradient\(to left,\s*var\(--gva-color-dark-panel-bg\),\s*transparent\);/,
  'dark mode should use the shared dark panel token for the table scroll fade'
)

console.log('statTableScrollFade tests passed')
