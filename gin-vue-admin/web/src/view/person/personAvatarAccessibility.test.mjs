import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./person.vue', import.meta.url), 'utf8')

const avatarUploadButton = source.match(/<button[^>]*class="update"[^>]*@click="openChooseImg"[^>]*>/)

assert.ok(avatarUploadButton, 'avatar upload control should use a native button element')
assert.match(avatarUploadButton[0], /type="button"/, 'avatar upload button should not submit forms')
assert.match(avatarUploadButton[0], /aria-label="重新上传头像"/, 'avatar upload button should have an explicit accessible name')
assert.match(avatarUploadButton[0], /@click="openChooseImg"/, 'avatar upload button should keep the image picker handler')

assert.doesNotMatch(source, /<span[^>]*class="update"[^>]*@click="openChooseImg"/, 'avatar upload control should not be a clickable span')

console.log('personAvatarAccessibility tests passed')
