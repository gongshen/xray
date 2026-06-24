import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync('src/view/system/state.vue', 'utf8')
const tableMatch = source.match(/<el-table[\s\S]*class="analysis-table"[\s\S]*?>/)

assert.ok(tableMatch, 'traffic analysis table should be rendered')
assert.match(tableMatch[0], /v-loading="trafficAnalysisLoading"/, 'traffic analysis table should show loading feedback')
assert.match(tableMatch[0], /:aria-busy="trafficAnalysisLoading"/, 'traffic analysis table should expose loading state to assistive tech')
assert.match(tableMatch[0], /aria-label="用户流量分钟明细表"/, 'traffic analysis table should have an explicit accessible name')

console.log('stateTrafficAnalysisAccessibility tests passed')
