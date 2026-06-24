import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./sysDictionaryDetail.vue', import.meta.url), 'utf8')

const labelInput = source.match(/<el-input[^>]*v-model="searchInfo\.label"[^>]*>/)
const valueInput = source.match(/<el-input[^>]*v-model="searchInfo\.value"[^>]*>/)
const statusSelect = source.match(/<el-select[^>]*v-model="searchInfo\.status"[^>]*>/)

assert.ok(labelInput, 'dictionary detail label search input should be rendered')
assert.match(labelInput[0], /aria-label="字典项展示值筛选"/, 'dictionary detail label search input should have an accessible name')

assert.ok(valueInput, 'dictionary detail value search input should be rendered')
assert.match(valueInput[0], /aria-label="字典项值筛选"/, 'dictionary detail value search input should have an accessible name')

assert.ok(statusSelect, 'dictionary detail status search select should be rendered')
assert.match(statusSelect[0], /aria-label="字典项启用状态筛选"/, 'dictionary detail status search select should have an accessible name')

console.log('sysDictionaryDetailAccessibility tests passed')
