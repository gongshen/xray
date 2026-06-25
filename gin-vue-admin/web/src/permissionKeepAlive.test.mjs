import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./permission.js', import.meta.url), 'utf8')

assert.match(source, /async function handleKeepAlive\(to\)/, 'permission guard should keep keep-alive preparation async')
assert.match(source, /to\.matched\?\.some\(item => item\.meta\?\.keepAlive\)/, 'keep-alive preparation should tolerate missing route metadata')
assert.match(source, /typeof element\?\.components\?\.default === 'function'/, 'keep-alive preparation should tolerate missing route components')
assert.match(source, /^\s*await handleKeepAlive\(to\)$/m, 'navigation guard should wait for keep-alive preparation before continuing')
assert.doesNotMatch(source, /^\s*handleKeepAlive\(to\)$/m, 'navigation guard should not fire keep-alive preparation without awaiting it')

console.log('permissionKeepAlive tests passed')
