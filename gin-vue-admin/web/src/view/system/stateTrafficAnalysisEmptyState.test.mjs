import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync('src/view/system/state.vue', 'utf8')
const table = source.match(/<el-table[\s\S]*class="analysis-table"[\s\S]*?>/)

assert.ok(table, 'traffic analysis table should be rendered')
assert.match(
  table[0],
  /:empty-text="trafficAnalysisLoading \? '加载中\.\.\.' : '暂无数据'"/,
  'traffic analysis table should not show an empty state while loading'
)

console.log('stateTrafficAnalysisEmptyState tests passed')
