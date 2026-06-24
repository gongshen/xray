import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./upload.vue', import.meta.url), 'utf8')

const editNameButton = source.match(/<button[^>]*class="name"[^>]*@click="editFileNameFunc\(scope\.row\)"[^>]*>/)

assert.ok(editNameButton, 'upload filename edit control should use a native button element')
assert.match(editNameButton[0], /type="button"/, 'upload filename edit button should not submit forms')
assert.match(editNameButton[0], /:aria-label="`编辑文件名或备注：\$\{scope\.row\.name\}`"/, 'upload filename edit button should include the current filename in its accessible name')
assert.match(editNameButton[0], /@click="editFileNameFunc\(scope\.row\)"/, 'upload filename edit button should keep the edit handler')

assert.doesNotMatch(source, /<div[^>]*class="name"[^>]*@click="editFileNameFunc\(scope\.row\)"/, 'upload filename edit control should not be a clickable div')

console.log('uploadAccessibility tests passed')
