import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync('src/view/system/state.vue', 'utf8')
const tableMatch = source.match(/<el-table[\s\S]*class="analysis-table"[\s\S]*?>/)

assert.ok(tableMatch, 'traffic analysis table should be rendered')
assert.match(tableMatch[0], /v-loading="trafficAnalysisLoading"/, 'traffic analysis table should show loading feedback')
assert.match(tableMatch[0], /:aria-busy="trafficAnalysisLoading"/, 'traffic analysis table should expose loading state to assistive tech')

console.log('stateTrafficAnalysisAccessibility tests passed')
