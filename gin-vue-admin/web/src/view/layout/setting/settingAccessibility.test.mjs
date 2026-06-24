import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./index.vue', import.meta.url), 'utf8')

const settingTrigger = source.match(/<el-button[^>]*class="drawer-container"[^>]*>/)
const themeGroup = source.match(/<div[^>]*class="theme-box"[^>]*>/)
const lightThemeButton = source.match(/<button[^>]*class="item"[^>]*@click="changeMode\('light'\)"[^>]*>/)
const darkThemeButton = source.match(/<button[^>]*class="item"[^>]*@click="changeMode\('dark'\)"[^>]*>/)

assert.ok(settingTrigger, 'setting drawer trigger should keep the existing Element Plus button')
assert.match(settingTrigger[0], /aria-label="打开系统配置"/, 'icon-only setting drawer trigger should have an accessible name')

assert.ok(themeGroup, 'theme mode controls should stay grouped visually')
assert.match(themeGroup[0], /role="group"/, 'theme mode controls should expose a semantic group')
assert.match(themeGroup[0], /aria-label="主题模式"/, 'theme mode group should have an accessible name')

assert.ok(lightThemeButton, 'light theme control should use a native button element')
assert.match(lightThemeButton[0], /type="button"/, 'light theme button should not submit forms')
assert.match(lightThemeButton[0], /:aria-pressed="userStore\.mode === 'light'"/, 'light theme button should expose selected state')
assert.match(lightThemeButton[0], /aria-label="切换为简约白主题"/, 'light theme button should have an explicit accessible name')

assert.ok(darkThemeButton, 'dark theme control should use a native button element')
assert.match(darkThemeButton[0], /type="button"/, 'dark theme button should not submit forms')
assert.match(darkThemeButton[0], /:aria-pressed="userStore\.mode === 'dark'"/, 'dark theme button should expose selected state')
assert.match(darkThemeButton[0], /aria-label="切换为商务黑主题"/, 'dark theme button should have an explicit accessible name')

assert.doesNotMatch(source, /<div[^>]*class="item"[^>]*@click="changeMode\('(light|dark)'\)"/, 'theme mode controls should not be clickable divs')

console.log('settingAccessibility tests passed')
