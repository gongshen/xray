import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./index.vue', import.meta.url), 'utf8')
const overlay = source.match(/<button[^>]*class="responsive-overlay"[^>]*>/)

assert.ok(overlay, 'responsive overlay should use a native button element')
assert.match(overlay[0], /type="button"/, 'responsive overlay button should not submit forms')
assert.match(overlay[0], /aria-label="关闭移动端菜单"/, 'responsive overlay should have an explicit close label')
assert.match(overlay[0], /@click="closeOverlay"/, 'responsive overlay should keep the close handler')
assert.doesNotMatch(source, /<div[^>]*class="responsive-overlay"[^>]*@click="closeOverlay"/, 'responsive overlay should not be a clickable div')

console.log('responsiveOverlayAccessibility tests passed')
