import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./binding.vue', import.meta.url), 'utf8')
const shareContainer = source.match(/<div[^>]*class="share-container"[^>]*>/)

assert.ok(shareContainer, 'binding share dialog should render a share content container')
assert.match(shareContainer[0], /v-loading="shareLoading"/, 'binding share dialog should show loading feedback')
assert.match(shareContainer[0], /:aria-busy="shareLoading"/, 'binding share dialog should expose loading state to assistive tech')
assert.match(shareContainer[0], /aria-label="分享配置内容"/, 'binding share dialog content should have an explicit accessible name')

console.log('bindingShareLoadingAccessibility tests passed')
