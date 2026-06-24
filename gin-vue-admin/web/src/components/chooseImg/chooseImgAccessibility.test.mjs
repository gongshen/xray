import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./index.vue', import.meta.url), 'utf8')

const editNameButton = source.match(/<button[^>]*class="img-title"[^>]*@click="editFileNameFunc\(item\)"[^>]*>/)

assert.ok(editNameButton, 'media library filename edit control should use a native button element')
assert.match(editNameButton[0], /type="button"/, 'media library filename edit button should not submit forms')
assert.match(editNameButton[0], /:aria-label="`编辑文件名或备注：\$\{item\.name\}`"/, 'media library filename edit button should include the current filename in its accessible name')
assert.match(editNameButton[0], /@click="editFileNameFunc\(item\)"/, 'media library filename edit button should keep the edit handler')

assert.doesNotMatch(source, /<div[^>]*class="img-title"[^>]*@click="editFileNameFunc\(item\)"/, 'media library filename edit control should not be a clickable div')

console.log('chooseImgAccessibility tests passed')
