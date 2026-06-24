import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./index.vue', import.meta.url), 'utf8')

const sidebar = source.match(/<el-aside[^>]*class="main-cont main-left gva-aside"[^>]*>/)
const collapseButton = source.match(/<button[^>]*class="menu-total"[^>]*>/)

assert.ok(sidebar, 'layout sidebar should keep its aside landmark')
assert.match(sidebar[0], /id="layout-sidebar"/, 'layout sidebar should expose a stable id for controls')

assert.ok(collapseButton, 'layout sidebar collapse control should use a native button element')
assert.match(collapseButton[0], /type="button"/, 'layout sidebar collapse button should not submit forms')
assert.match(collapseButton[0], /aria-controls="layout-sidebar"/, 'layout sidebar collapse button should reference the controlled sidebar')
assert.match(collapseButton[0], /:aria-expanded="!isCollapse"/, 'layout sidebar collapse button should expose sidebar expansion state')
assert.match(collapseButton[0], /:aria-label="isCollapse \? '展开侧边栏' : '收起侧边栏'"/, 'layout sidebar collapse button should have a state-aware accessible name')
assert.match(collapseButton[0], /@click="totalCollapse"/, 'layout sidebar collapse button should keep the collapse handler')

assert.doesNotMatch(source, /<div[^>]*class="menu-total"[^>]*@click="totalCollapse"/, 'layout sidebar collapse control should not be a clickable div')

console.log('layoutAccessibility tests passed')
