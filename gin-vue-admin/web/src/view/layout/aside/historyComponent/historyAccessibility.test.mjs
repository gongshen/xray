import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./history.vue', import.meta.url), 'utf8')

const contextMenu = source.match(/<ul[^>]*class="contextmenu"[^>]*>/)
const expectations = [
  { handler: 'closeAll', label: '关闭所有标签' },
  { handler: 'closeLeft', label: '关闭左侧标签' },
  { handler: 'closeRight', label: '关闭右侧标签' },
  { handler: 'closeOther', label: '关闭其他标签' },
]

assert.ok(contextMenu, 'history context menu should be rendered')
assert.match(contextMenu[0], /role="menu"/, 'history context menu should expose menu semantics')
assert.match(contextMenu[0], /aria-label="标签页操作"/, 'history context menu should have an accessible name')

for (const { handler, label } of expectations) {
  const button = source.match(new RegExp(`<button[^>]*class="contextmenu-action"[^>]*@click="${handler}"[^>]*>`))
  assert.ok(button, `${label} action should use a native button`)
  assert.match(button[0], /type="button"/, `${label} action should not submit forms`)
  assert.match(button[0], /role="menuitem"/, `${label} action should expose menuitem semantics`)
}

assert.doesNotMatch(source, /<li[^>]*@click="(closeAll|closeLeft|closeRight|closeOther)"/, 'history context menu actions should not be clickable list items')

console.log('historyAccessibility tests passed')
