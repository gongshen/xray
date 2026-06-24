import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./binding.vue', import.meta.url), 'utf8')
const refreshIcons = [...source.matchAll(/<el-icon[^>]*class="auto-icon"[^>]*>/g)].map((match) => match[0])

assert.equal(refreshIcons.length, 2, 'binding filters should render two refresh icons')

const expectations = [
  { clickHandler: 'getSrvs', label: '刷新服务器选项' },
  { clickHandler: 'getUsers', label: '刷新用户选项' },
]

for (const { clickHandler, label } of expectations) {
  const icon = refreshIcons.find((candidate) => candidate.includes(`@click="${clickHandler}"`))
  assert.ok(icon, `${label} icon should be rendered`)
  assert.match(icon, /role="button"/, `${label} icon should expose button semantics`)
  assert.match(icon, /tabindex="0"/, `${label} icon should be keyboard reachable`)
  assert.match(icon, new RegExp(`aria-label="${label}"`), `${label} icon should have an accessible name`)
  assert.match(icon, new RegExp(`@keydown\\.enter\\.prevent="${clickHandler}"`), `${label} icon should support Enter`)
  assert.match(icon, new RegExp(`@keydown\\.space\\.prevent="${clickHandler}"`), `${label} icon should support Space`)
}

console.log('bindingRefreshIconAccessibility tests passed')
