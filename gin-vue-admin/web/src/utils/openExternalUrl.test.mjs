import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./openExternalUrl.js', import.meta.url), 'utf8')
const module = { exports: {} }
const opened = []
const window = {
  open: (...args) => {
    opened.push(args)
  }
}

Function('module', 'window', source.replace(/export const /g, 'const ') + '; module.exports = { normalizeExternalUrl, isExternalHttpUrl, openExternalUrl }')(module, window)

const { normalizeExternalUrl, isExternalHttpUrl, openExternalUrl } = module.exports

assert.equal(normalizeExternalUrl('  https://example.com/path  '), 'https://example.com/path')
assert.equal(normalizeExternalUrl(null), '')
assert.equal(isExternalHttpUrl('https://example.com'), true)
assert.equal(isExternalHttpUrl('http://example.com'), true)
assert.equal(isExternalHttpUrl('javascript:alert(1)'), false)
assert.equal(isExternalHttpUrl('/internal/path'), false)
assert.equal(isExternalHttpUrl(''), false)

assert.equal(openExternalUrl(' https://example.com/docs '), true)
assert.deepEqual(opened, [['https://example.com/docs', '_blank', 'noopener,noreferrer']])
assert.equal(openExternalUrl('javascript:alert(1)'), false)
assert.equal(openExternalUrl('/dashboard'), false)
assert.equal(opened.length, 1, 'unsafe or internal URLs should not open a new window')

console.log('openExternalUrl tests passed')
