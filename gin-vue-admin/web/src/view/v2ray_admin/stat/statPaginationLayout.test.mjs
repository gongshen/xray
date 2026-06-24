import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync('src/view/v2ray_admin/stat/statPage.scss', 'utf8')

assert.match(
  source,
  /:deep\(\.el-pagination\)\s*\{[\s\S]*flex-wrap:\s*wrap;/,
  'pagination controls should wrap instead of overflowing narrow containers'
)
assert.match(
  source,
  /:deep\(\.el-pagination\)\s*\{[\s\S]*justify-content:\s*center;/,
  'pagination controls should remain centered when they wrap'
)
assert.match(
  source,
  /:deep\(\.el-pagination\)\s*\{[\s\S]*gap:\s*8px;/,
  'wrapped pagination controls should keep stable spacing'
)

console.log('statPaginationLayout tests passed')
