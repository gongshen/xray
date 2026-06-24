import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./user.vue', import.meta.url), 'utf8')

const headerImageButton = source.match(/<button[^>]*class="header-img-trigger"[^>]*@click="openHeaderChange"[^>]*>/)

assert.ok(headerImageButton, 'user header image chooser should use a native button element')
assert.match(headerImageButton[0], /type="button"/, 'user header image chooser should not submit forms')
assert.match(headerImageButton[0], /aria-label="选择用户头像"/, 'user header image chooser should have an explicit accessible name')
assert.match(headerImageButton[0], /@click="openHeaderChange"/, 'user header image chooser should keep the media picker handler')

assert.doesNotMatch(source, /<div[^>]*@click="openHeaderChange"/, 'user header image chooser should not be a clickable div')

console.log('userHeaderImageAccessibility tests passed')
