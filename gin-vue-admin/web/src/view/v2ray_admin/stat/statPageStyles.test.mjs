import assert from 'node:assert/strict'
import fs from 'node:fs'

const sharedStylePath = 'src/view/v2ray_admin/stat/statPage.scss'
const adminStatPath = 'src/view/v2ray_admin/stat/stat.vue'
const userStatPath = 'src/view/v2ray/stat/stat.vue'

assert.equal(fs.existsSync(sharedStylePath), true, 'shared statPage.scss should exist')

const sharedStyle = fs.readFileSync(sharedStylePath, 'utf8')
const requiredSelectors = [
  '.page',
  '.charts-section',
  '.chart-card',
  '.gva-table-box',
  '.traffic-down',
  '.traffic-up',
  '@media (prefers-reduced-motion: reduce)',
  '@media screen and (max-width: 768px)',
]

for (const selector of requiredSelectors) {
  assert.match(sharedStyle, new RegExp(selector.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')), `missing ${selector}`)
}

const expectedStyleRefs = [
  [adminStatPath, '<style scoped lang="scss" src="./statPage.scss"></style>'],
  [userStatPath, '<style scoped lang="scss" src="../../v2ray_admin/stat/statPage.scss"></style>'],
]

for (const [file, expectedStyleRef] of expectedStyleRefs) {
  const source = fs.readFileSync(file, 'utf8')
  assert.match(source, new RegExp(expectedStyleRef.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')), `${file} should use shared stat page styles`)
  assert.doesNotMatch(source, /<style scoped>[\s\S]*\.chart-card[\s\S]*<\/style>/, `${file} should not duplicate chart styles inline`)
  assert.doesNotMatch(source, /<style scoped>[\s\S]*\.gva-table-box[\s\S]*<\/style>/, `${file} should not duplicate table styles inline`)
}

console.log('statPageStyles tests passed')
