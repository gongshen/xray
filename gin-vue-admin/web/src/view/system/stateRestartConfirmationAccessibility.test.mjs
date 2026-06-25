import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./state.vue', import.meta.url), 'utf8')
const text = (...codes) => String.fromCharCode(...codes)
const escaped = (value) => Array.from(value).map((char) => {
  return '\\u' + char.charCodeAt(0).toString(16).padStart(4, '0')
}).join('')
const includesLiteralString = (haystack, value) => {
  return haystack.includes(`'${value}'`) || haystack.includes(`'${escaped(value)}'`)
}
const includesPropertyString = (haystack, property, value) => {
  return haystack.includes(`${property}: '${value}'`) || haystack.includes(`${property}: '${escaped(value)}'`)
}

const restartButtonLabel = text(0x91cd, 0x542f, 0x5f53, 0x524d, 0x670d, 0x52a1, 0x5668)
const restartMessage = text(0x786e, 0x5b9a, 0x8981, 0x91cd, 0x542f, 0x5f53, 0x524d, 0x670d, 0x52a1, 0x5668, 0x5417, 0xff1f, 0x91cd, 0x542f, 0x53ef, 0x80fd, 0x9700, 0x8981, 0x51e0, 0x5206, 0x949f, 0x65f6, 0x95f4, 0x3002)
const restartTitle = text(0x91cd, 0x542f, 0x670d, 0x52a1, 0x5668)
const restartAction = text(0x91cd, 0x542f)
const cancelAction = text(0x53d6, 0x6d88)

const restartButton = source.match(/<el-button[^>]*@click="restartVPS"[^>]*>/)
const restartConfirmation = source.match(/ElMessageBox\.confirm\([\s\S]*?restartVPSApi/)

assert.ok(restartButton, 'server restart button should be rendered')
assert.match(restartButton[0], /type="danger"/, 'server restart button should use danger styling')
assert.ok(restartButton[0].includes(`aria-label="${restartButtonLabel}"`), 'server restart button should have a specific accessible name')

assert.ok(restartConfirmation, 'server restart should ask for confirmation')
assert.ok(includesLiteralString(restartConfirmation[0], restartMessage), 'server restart confirmation should name the target and duration risk')
assert.ok(includesLiteralString(restartConfirmation[0], restartTitle), 'server restart confirmation should use a specific title')
assert.ok(includesPropertyString(restartConfirmation[0], 'confirmButtonText', restartAction), 'server restart confirmation should use action-specific text')
assert.ok(includesPropertyString(restartConfirmation[0], 'cancelButtonText', cancelAction), 'server restart confirmation should keep a cancel action')
assert.match(restartConfirmation[0], /type:\s*'warning'/, 'server restart confirmation should keep warning intent')

console.log('stateRestartConfirmationAccessibility tests passed')
