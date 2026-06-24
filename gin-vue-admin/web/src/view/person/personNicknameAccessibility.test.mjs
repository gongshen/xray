import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./person.vue', import.meta.url), 'utf8')

const editNickButton = source.match(/<button[^>]*class="pointer nick-action"[^>]*@click="openEdit"[^>]*>/)
const saveNickButton = source.match(/<button[^>]*class="pointer nick-action"[^>]*@click="enterEdit"[^>]*>/)
const cancelNickButton = source.match(/<button[^>]*class="pointer nick-action"[^>]*@click="closeEdit"[^>]*>/)

assert.ok(editNickButton, 'nickname edit control should use a native button element')
assert.match(editNickButton[0], /type="button"/, 'nickname edit button should not submit forms')
assert.match(editNickButton[0], /aria-label="编辑昵称"/, 'nickname edit button should have an explicit accessible name')

assert.ok(saveNickButton, 'nickname save control should use a native button element')
assert.match(saveNickButton[0], /type="button"/, 'nickname save button should not submit forms')
assert.match(saveNickButton[0], /aria-label="保存昵称"/, 'nickname save button should have an explicit accessible name')

assert.ok(cancelNickButton, 'nickname cancel control should use a native button element')
assert.match(cancelNickButton[0], /type="button"/, 'nickname cancel button should not submit forms')
assert.match(cancelNickButton[0], /aria-label="取消编辑昵称"/, 'nickname cancel button should have an explicit accessible name')

assert.doesNotMatch(source, /<el-icon[^>]*class="pointer"[^>]*@click="(openEdit|enterEdit|closeEdit)"/, 'nickname actions should not be clickable icon components')

console.log('personNicknameAccessibility tests passed')
