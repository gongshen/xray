import assert from 'node:assert/strict'
import fs from 'node:fs'

const files = [
  'src/view/v2ray_admin/stat/stat.vue',
  'src/view/v2ray/stat/stat.vue',
]

for (const file of files) {
  const source = fs.readFileSync(file, 'utf8')
  const startPicker = source.match(/<el-date-picker[^>]*v-model="searchInfo\.startCreatedAt"[^>]*>/)
  const endPicker = source.match(/<el-date-picker[^>]*v-model="searchInfo\.endCreatedAt"[^>]*>/)

  assert.ok(startPicker, `${file} should render a start date picker`)
  assert.ok(endPicker, `${file} should render an end date picker`)
  assert.match(startPicker[0], /aria-label="统计开始日期"/, `${file} start date picker should have an accessible name`)
  assert.match(endPicker[0], /aria-label="统计结束日期"/, `${file} end date picker should have an accessible name`)
}

console.log('statDatePickerAccessibility tests passed')
