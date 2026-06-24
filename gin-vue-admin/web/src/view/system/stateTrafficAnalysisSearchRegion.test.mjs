import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync('src/view/system/state.vue', 'utf8')
const form = source
  .split('\n')
  .find(line => line.includes('<el-form') && line.includes('class="analysis-form"'))

assert.ok(form, 'traffic analysis filter form should be rendered')
assert.match(form, /role="search"/, 'traffic analysis filter form should expose search semantics')
assert.match(form, /aria-label="流量分析筛选"/, 'traffic analysis filter form should have an explicit search landmark name')

console.log('stateTrafficAnalysisSearchRegion tests passed')
