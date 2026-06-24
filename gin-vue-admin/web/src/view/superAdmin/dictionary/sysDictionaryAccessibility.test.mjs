import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./sysDictionary.vue', import.meta.url), 'utf8')

const nameInput = source.match(/<el-input[^>]*v-model="searchInfo\.name"[^>]*>/)
const typeInput = source.match(/<el-input[^>]*v-model="searchInfo\.type"[^>]*>/)
const statusSelect = source.match(/<el-select[^>]*v-model="searchInfo\.status"[^>]*>/)
const descInput = source.match(/<el-input[^>]*v-model="searchInfo\.desc"[^>]*>/)

assert.ok(nameInput, 'dictionary name search input should be rendered')
assert.match(nameInput[0], /aria-label="字典中文名筛选"/, 'dictionary name search input should have an accessible name')

assert.ok(typeInput, 'dictionary type search input should be rendered')
assert.match(typeInput[0], /aria-label="字典英文名筛选"/, 'dictionary type search input should have an accessible name')

assert.ok(statusSelect, 'dictionary status search select should be rendered')
assert.match(statusSelect[0], /aria-label="字典状态筛选"/, 'dictionary status search select should have an accessible name')

assert.ok(descInput, 'dictionary description search input should be rendered')
assert.match(descInput[0], /aria-label="字典描述筛选"/, 'dictionary description search input should have an accessible name')

console.log('sysDictionaryAccessibility tests passed')
