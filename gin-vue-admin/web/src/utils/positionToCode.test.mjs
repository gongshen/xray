import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./positionToCode.js', import.meta.url), 'utf8')

assert.match(source, /let disposeDomMouseDown = null/, 'position-to-code should track the active document listener')
assert.match(source, /document\.addEventListener\('mousedown', handleMouseDown\)/, 'position-to-code should register a mousedown listener without replacing document handlers')
assert.match(source, /document\.removeEventListener\('mousedown', handleMouseDown\)/, 'position-to-code should expose listener cleanup')
assert.match(source, /export const disposeDom = \(\) =>/, 'position-to-code should export cleanup for the development listener')
assert.doesNotMatch(source, /document\.onmousedown/, 'position-to-code should not overwrite document.onmousedown')
assert.match(source, /if \(!filePath\) \{/, 'position-to-code should ignore missing file paths')
assert.match(source, /encodeURIComponent\(filePath\)/, 'position-to-code should encode file paths in the editor-open request')

console.log('positionToCode tests passed')
