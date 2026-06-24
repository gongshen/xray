import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./index.vue', import.meta.url), 'utf8')

const dashboardIllustration = source.match(/<img[^>]*src="@\/assets\/dashboard\.png"[^>]*>/)

assert.ok(dashboardIllustration, 'dashboard should render the top card illustration')
assert.match(dashboardIllustration[0], /alt="管理后台工作台插图"/, 'dashboard top card illustration should have descriptive alt text')

console.log('dashboardAccessibility tests passed')
