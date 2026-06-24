import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync('src/view/system/state.vue', 'utf8')

const controls = [
  ['user_tag', 'el-select', 'trafficAnalysisForm.user_tag', '流量分析用户'],
  ['date', 'el-date-picker', 'trafficAnalysisForm.date', '流量分析日期'],
  ['start', 'el-input', 'trafficAnalysisForm.start', '流量分析开始时间'],
  ['end', 'el-input', 'trafficAnalysisForm.end', '流量分析结束时间'],
]

for (const [name, tag, model, label] of controls) {
  const modelPattern = model.replaceAll('.', '\\.')
  const control = source.match(new RegExp(`<${tag}\\b(?=[^>]*v-model="${modelPattern}")[^>]*>`))
  assert.ok(control, `traffic analysis ${name} control should be rendered`)
  assert.match(control[0], new RegExp(`aria-label="${label}"`), `${label} control should have an explicit accessible name`)
}

const actionGroup = source.match(/<el-form-item[^>]*class="analysis-action"[^>]*>/)
assert.ok(actionGroup, 'traffic analysis query action group should be rendered')
assert.match(actionGroup[0], /role="group"/, 'traffic analysis query actions should expose group semantics')
assert.match(actionGroup[0], /aria-label="流量分析查询操作"/, 'traffic analysis query actions should have an explicit group name')

console.log('stateTrafficAnalysisFormAccessibility tests passed')
