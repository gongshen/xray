import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync('src/view/system/state.vue', 'utf8')

const progressByMetric = [
  ['diskPercent', '磁盘使用率'],
  ['cpuPercent', 'CPU 使用率'],
  ['memPercent', '内存使用率'],
]

for (const [model, label] of progressByMetric) {
  const progress = source.match(new RegExp(`<el-progress[\\s\\S]*:percentage="${model}"[\\s\\S]*?>`))
  assert.ok(progress, `${label} progress should be rendered`)
  assert.match(progress[0], new RegExp(`aria-label="${label}"`), `${label} progress should have an explicit accessible name`)
}

console.log('stateMetricProgressAccessibility tests passed')
