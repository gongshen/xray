import assert from 'node:assert/strict'
import fs from 'node:fs'

const files = [
  'src/view/v2ray_admin/stat/stat.vue',
  'src/view/v2ray/stat/stat.vue',
]

for (const file of files) {
  const source = fs.readFileSync(file, 'utf8')
  const actionGroup = source.match(/<el-form-item\b[^>]*class="stat-action-group"[^>]*>/)

  assert.ok(actionGroup, `${file} should render a grouped stat action area`)
  assert.match(actionGroup[0], /role="group"/, `${file} stat actions should expose group semantics`)
  assert.match(actionGroup[0], /aria-label="统计查询操作"/, `${file} stat actions should have an explicit group name`)
}

console.log('statActionGroupAccessibility tests passed')
