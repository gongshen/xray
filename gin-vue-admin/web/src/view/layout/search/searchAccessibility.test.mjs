import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./search.vue', import.meta.url), 'utf8')

const refreshButton = source.match(/<button[^>]*class="gvaIcon gvaIcon-refresh[^"]*"[^>]*>/)
const searchButton = source.match(/<button[^>]*class="gvaIcon gvaIcon-search[^"]*"[^>]*>/)

assert.ok(refreshButton, 'layout refresh control should use a native button element')
assert.match(refreshButton[0], /type="button"/, 'layout refresh button should not submit forms')
assert.match(refreshButton[0], /aria-label="刷新当前页面"/, 'layout refresh button should have an explicit accessible name')
assert.match(refreshButton[0], /@click="handleReload"/, 'layout refresh button should keep the reload handler')
assert.match(refreshButton[0], /:aria-busy="reload"/, 'layout refresh button should expose reload state')

assert.ok(searchButton, 'layout search control should use a native button element')
assert.match(searchButton[0], /type="button"/, 'layout search button should not submit forms')
assert.match(searchButton[0], /aria-label="打开页面搜索"/, 'layout search button should have an explicit accessible name')
assert.match(searchButton[0], /@click="showSearch"/, 'layout search button should keep the search handler')

assert.doesNotMatch(source, /<div[^>]*class="gvaIcon gvaIcon-refresh"[^>]*@click="handleReload"/, 'layout refresh control should not be a clickable div')
assert.doesNotMatch(source, /<div[^>]*class="gvaIcon gvaIcon-search"[^>]*@click="showSearch"/, 'layout search control should not be a clickable div')

console.log('searchAccessibility tests passed')
