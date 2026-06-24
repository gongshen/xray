import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./index.vue', import.meta.url), 'utf8')
const toggle = source.match(/<button[^>]*class="mobile-menu-toggle"[^>]*>/)

assert.ok(toggle, 'mobile menu toggle should use a native button element')
assert.match(toggle[0], /type="button"/, 'mobile menu toggle button should not submit forms')
assert.match(toggle[0], /@click="toggleMobileMenu"/, 'mobile menu toggle should keep the click handler')
assert.match(toggle[0], /:aria-expanded="isMenuVisible"/, 'mobile menu toggle should expose expanded state')
assert.match(toggle[0], /:aria-label="isMenuVisible \? '关闭移动端菜单' : '打开移动端菜单'"/, 'mobile menu toggle should expose an action-specific accessible name')
assert.doesNotMatch(source, /<div[^>]*class="mobile-menu-toggle"[^>]*@click="toggleMobileMenu"/, 'mobile menu toggle should not be a clickable div')

console.log('mobileMenuToggleAccessibility tests passed')
