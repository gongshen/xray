import assert from 'node:assert/strict'
import fs from 'node:fs'

const sharedStyle = fs.readFileSync('src/view/v2ray_admin/stat/statPage.scss', 'utf8')

assert.doesNotMatch(sharedStyle, /transition:\s*all\b/, 'chart cards should not transition every animatable property')
assert.match(
  sharedStyle,
  /\.chart-card\s*\{[\s\S]*transition:\s*box-shadow\s+0\.2s\s+ease;/,
  'chart cards should only animate their hover shadow'
)
assert.match(
  sharedStyle,
  /@media \(prefers-reduced-motion:\s*reduce\)[\s\S]*\.chart-card,[\s\S]*\.gva-table-box\s*\{[\s\S]*transition:\s*none;/,
  'reduced motion should disable stat page transitions'
)

console.log('statMotionStyles tests passed')
