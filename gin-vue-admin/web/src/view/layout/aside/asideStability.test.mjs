import assert from 'node:assert/strict'
import fs from 'node:fs'

const asideSource = fs.readFileSync(new URL('./index.vue', import.meta.url), 'utf8')
const asideComponentSource = fs.readFileSync(new URL('./asideComponent/index.vue', import.meta.url), 'utf8')

assert.match(asideSource, /v-for="item in menuRouters"/, 'aside should render from a guarded menu list')
assert.match(asideSource, /<template\b[^>]*v-for="item in menuRouters"[^>]*:key="item.name"[^>]*>/s, 'aside template fragment should have a stable key')
assert.match(asideSource, /const menuRouters = computed\(\(\) => routerStore\.asyncRouters\[0\]\?\.children \|\| \[\]\)/, 'aside should tolerate async routes before they are loaded')
assert.doesNotMatch(asideSource, /routerStore\.asyncRouters\[0\]\.children/, 'aside template should not directly index async route children')

assert.match(asideComponentSource, /default: \(\) => \(\{\}\)/, 'aside component should default routerInfo to an object')
assert.match(asideComponentSource, /const visibleChildren = computed\(\(\) => props\.routerInfo\.children\?\.filter\(item => !item\.hidden\) \|\| \[\]\)/, 'aside component should derive visible children safely')
assert.match(asideComponentSource, /v-for="item in visibleChildren"/, 'aside component should recurse only over visible guarded children')

console.log('asideStability tests passed')
