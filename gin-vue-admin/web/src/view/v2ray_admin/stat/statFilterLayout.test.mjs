import assert from 'node:assert/strict'
import fs from 'node:fs'

const adminStat = fs.readFileSync('src/view/v2ray_admin/stat/stat.vue', 'utf8')
assert.doesNotMatch(adminStat, /style="width:\s*194px"/, 'stat filters should not keep fixed widths inline')
assert.equal(
  (adminStat.match(/class="stat-filter-select"/g) || []).length,
  2,
  'admin stat user/server filters should share stat-filter-select class'
)

const userStat = fs.readFileSync('src/view/v2ray/stat/stat.vue', 'utf8')
assert.doesNotMatch(userStat, /style="width:\s*194px"/, 'user stat should not introduce fixed inline filter widths')

const sharedStyle = fs.readFileSync('src/view/v2ray_admin/stat/statPage.scss', 'utf8')
assert.match(sharedStyle, /\.stat-filter-select\s*\{[\s\S]*width:\s*194px;/, 'shared styles should own desktop filter width')
assert.match(sharedStyle, /@media screen and \(max-width:\s*768px\)[\s\S]*\.stat-filter-select\s*\{[\s\S]*width:\s*100%;/, 'mobile filters should stretch to container width')

console.log('statFilterLayout tests passed')
