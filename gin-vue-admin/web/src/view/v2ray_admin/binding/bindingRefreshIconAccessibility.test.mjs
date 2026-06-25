import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./binding.vue', import.meta.url), 'utf8')
const refreshButtons = [...source.matchAll(/<button[^>]*class="auto-icon"[^>]*>[\s\S]*?<\/button>/g)].map((match) => match[0])

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
  assert.match(button, /<el-icon>\s*<refresh\s*\/>\s*<\/el-icon>/, `${label} button should wrap the refresh svg in el-icon sizing`)
}

assert.doesNotMatch(source, /<el-icon[^>]*class="auto-icon"[^>]*@click="(getSrvs|getUsers)"/, 'binding refresh actions should not be clickable icon components')
assert.match(source, /\.auto-icon\s+:deep\(svg\)\s*\{[\s\S]*width:\s*1em;[\s\S]*height:\s*1em;/, 'binding refresh svg should be constrained to text-sized dimensions')

console.log('bindingRefreshIconAccessibility tests passed')
