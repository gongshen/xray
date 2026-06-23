import assert from 'node:assert/strict'
import fs from 'node:fs'

const statFiles = [
  'src/view/v2ray_admin/stat/stat.vue',
  'src/view/v2ray/stat/stat.vue',
]

for (const file of statFiles) {
  const source = fs.readFileSync(file, 'utf8')
  assert.match(source, /<el-table[\s\S]*class="stat-table"/, `${file} should mark the stat table with shared class`)
  assert.doesNotMatch(source, /<el-table[\s\S]*style="width:\s*100%;\s*min-height:\s*200px;"/, `${file} should not keep stat table sizing inline`)
}

const sharedStyle = fs.readFileSync('src/view/v2ray_admin/stat/statPage.scss', 'utf8')
assert.match(sharedStyle, /\.stat-table\s*\{[\s\S]*width:\s*100%;[\s\S]*min-height:\s*200px;/, 'statPage.scss should own table sizing')
assert.match(sharedStyle, /@media screen and \(max-width:\s*768px\)[\s\S]*overflow-x:\s*auto;/, 'mobile table container should scroll horizontally')
assert.match(sharedStyle, /@media screen and \(max-width:\s*768px\)[\s\S]*\.stat-table\s*\{[\s\S]*min-width:\s*680px;/, 'mobile stat table should keep readable columns')

console.log('statTableLayout tests passed')
