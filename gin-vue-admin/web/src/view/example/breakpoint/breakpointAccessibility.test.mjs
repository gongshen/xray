import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./breakpoint.vue', import.meta.url), 'utf8')

const fileUploadButton = source.match(/<button[^>]*class="fileUpload"[^>]*@click="inputChange"[^>]*>/)

assert.ok(fileUploadButton, 'breakpoint file chooser should use a native button element')
assert.match(fileUploadButton[0], /type="button"/, 'breakpoint file chooser button should not submit forms')
assert.match(fileUploadButton[0], /aria-label="选择断点续传文件"/, 'breakpoint file chooser should have an explicit accessible name')
assert.match(fileUploadButton[0], /@click="inputChange"/, 'breakpoint file chooser button should keep the input trigger handler')

assert.doesNotMatch(source, /<div[^>]*class="fileUpload"[^>]*@click="inputChange"/, 'breakpoint file chooser should not be a clickable div')

console.log('breakpointAccessibility tests passed')
