import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./index.vue', import.meta.url), 'utf8')
const toggle = source.match(/<button[^>]*class="screenfull-toggle"[^>]*>/)

assert.ok(toggle, 'screenfull toggle should use a native button element')
assert.match(toggle[0], /type="button"/, 'screenfull toggle button should not submit forms')
assert.match(toggle[0], /@click="clickFull"/, 'screenfull toggle should keep the click handler')
assert.match(toggle[0], /:aria-label="isShow \? '进入全屏' : '退出全屏'"/, 'screenfull toggle should expose action-specific accessible names')
assert.doesNotMatch(source, /<div[^>]*@click="clickFull"/, 'screenfull toggle should not be a clickable div')

console.log('screenfullAccessibility tests passed')
