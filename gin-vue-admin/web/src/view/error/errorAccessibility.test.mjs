import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./index.vue', import.meta.url), 'utf8')

const notFoundImage = source.match(/<img[^>]*src="\.\.\/\.\.\/assets\/notFound\.png"[^>]*>/)
const guideImage = source.match(/<img[^>]*src="\.\.\/\.\.\/assets\/qm\.png"[^>]*>/)

assert.ok(notFoundImage, 'error page should render the not-found illustration')
assert.match(notFoundImage[0], /alt="页面未找到插图"/, 'not-found illustration should have descriptive alt text')

assert.ok(guideImage, 'error page should render the route guidance illustration')
assert.match(guideImage[0], /alt="角色权限配置引导图"/, 'route guidance illustration should have descriptive alt text')

console.log('errorAccessibility tests passed')
