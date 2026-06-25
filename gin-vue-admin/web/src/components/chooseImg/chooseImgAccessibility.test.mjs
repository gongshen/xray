import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./index.vue', import.meta.url), 'utf8')

const editNameButton = source.match(/<button[^>]*class="img-title"[^>]*@click="editFileNameFunc\(item\)"[^>]*>/)
const chooseImageButton = source.match(/<button[^>]*class="choose-img-trigger"[^>]*@click="chooseImg\(item\.url,target,targetKey\)"[^>]*>/)
const searchInput = source.match(/<el-input[^>]*v-model="search\.keyword"[^>]*>/)
const chooseImageStyle = source.match(/\.choose-img-trigger\s*\{([\s\S]*?)\n\s*\}/)

assert.ok(chooseImageButton, 'media library image choose control should use a native button element')
assert.match(chooseImageButton[0], /type="button"/, 'media library image choose button should not submit forms')
assert.match(chooseImageButton[0], /:aria-label="`选择图片：\$\{item\.name\}`"/, 'media library image choose button should include the current filename in its accessible name')
assert.doesNotMatch(source, /<el-image[^>]*@click="chooseImg\(item\.url,target,targetKey\)"/, 'media library image should not be directly clickable')

assert.ok(chooseImageStyle, 'media library image choose button should have overlay styles')
for (const declaration of [
  'position: absolute;',
  'inset: 0;',
  'border: 0;',
  'background: transparent;',
  'padding: 0;',
  'cursor: pointer;'
]) {
  assert.ok(
    chooseImageStyle[1].includes(declaration),
    `media library image choose button style should include ${declaration}`
  )
}

assert.ok(editNameButton, 'media library filename edit control should use a native button element')
assert.match(editNameButton[0], /type="button"/, 'media library filename edit button should not submit forms')
assert.match(editNameButton[0], /:aria-label="`编辑文件名或备注：\$\{item\.name\}`"/, 'media library filename edit button should include the current filename in its accessible name')
assert.match(editNameButton[0], /@click="editFileNameFunc\(item\)"/, 'media library filename edit button should keep the edit handler')

assert.doesNotMatch(source, /<div[^>]*class="img-title"[^>]*@click="editFileNameFunc\(item\)"/, 'media library filename edit control should not be a clickable div')

assert.ok(searchInput, 'media library search input should be rendered')
assert.match(searchInput[0], /aria-label="媒体库文件搜索"/, 'media library search input should have an accessible name')

console.log('chooseImgAccessibility tests passed')
