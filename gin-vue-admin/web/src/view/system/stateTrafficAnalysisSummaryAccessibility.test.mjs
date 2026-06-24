import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync('src/view/system/state.vue', 'utf8')
const summary = source.match(/<div[^>]*v-if="trafficAnalysisResult"[^>]*class="analysis-summary"[^>]*>/)

assert.ok(summary, 'traffic analysis summary should be rendered when results exist')
assert.match(summary[0], /role="status"/, 'traffic analysis summary should expose status semantics')
assert.match(summary[0], /aria-live="polite"/, 'traffic analysis summary should announce updates politely')
assert.match(summary[0], /aria-label="流量分析摘要"/, 'traffic analysis summary should have an explicit accessible name')

console.log('stateTrafficAnalysisSummaryAccessibility tests passed')
