import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./binding.vue', import.meta.url), 'utf8')
const refreshButtons = [...source.matchAll(/<button[^>]*class="auto-icon"[^>]*>/g)].map((match) => match[0])

assert.equal(refreshButtons.length, 2, 'binding filters should render two refresh buttons')

const expectations = [
  { clickHandler: 'getSrvs', label: '刷新服务器选项' },
  { clickHandler: 'getUsers', label: '刷新用户选项' },
]

for (const { clickHandler, label } of expectations) {
  const button = refreshButtons.find((candidate) => candidate.includes(`@click="${clickHandler}"`))
  assert.ok(button, `${label} button should be rendered`)
  assert.match(button, /type="button"/, `${label} button should not submit forms`)
  assert.match(button, new RegExp(`aria-label="${label}"`), `${label} button should have an accessible name`)
  assert.match(button, new RegExp(`@click="${clickHandler}"`), `${label} button should keep the refresh handler`)
}

assert.doesNotMatch(source, /<el-icon[^>]*class="auto-icon"[^>]*@click="(getSrvs|getUsers)"/, 'binding refresh actions should not be clickable icon components')

console.log('bindingRefreshIconAccessibility tests passed')
