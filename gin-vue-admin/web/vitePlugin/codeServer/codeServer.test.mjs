import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./index.js', import.meta.url), 'utf8')
  .replace("import { spawn } from 'node:child_process'", 'const spawn = () => ({ on() {}, unref() {} })')
  .replace('export default function GvaPositionServer', 'function GvaPositionServer')
  .replace(/export const /g, 'const ')
const module = { exports: {} }
Function('module', source + '; module.exports = { getEditorLaunch, parseEditorTarget }')(module)

const { getEditorLaunch, parseEditorTarget } = module.exports

assert.equal(parseEditorTarget(''), null)
assert.equal(parseEditorTarget('null'), null)
assert.equal(parseEditorTarget('D:/repo/src/App.vue'), null)
assert.equal(parseEditorTarget('D:/repo/src/App.vue:abc'), null)
assert.equal(parseEditorTarget('D:/repo/src/App.vue:12&calc'), null)

const target = parseEditorTarget('D:/repo/src/App.vue:12')
assert.deepEqual(target, {
  filePath: 'D:/repo/src/App.vue',
  line: '12',
  editorTarget: 'D:/repo/src/App.vue:12',
})

assert.deepEqual(getEditorLaunch(target, 'webstorm', 'win32'), {
  command: 'webstorm64.exe',
  args: ['--line', '12', 'D:/repo/src/App.vue'],
  shell: false,
})

assert.deepEqual(getEditorLaunch(target, 'code', 'linux'), {
  command: 'code',
  args: ['-r', '-g', 'D:/repo/src/App.vue:12'],
  shell: false,
})

assert.deepEqual(getEditorLaunch(target, 'code', 'win32'), {
  command: 'cmd.exe',
  args: ['/d', '/s', '/c', 'code', '-r', '-g', 'D:/repo/src/App.vue:12'],
  shell: false,
})

assert.equal(getEditorLaunch(null), null)

console.log('codeServer tests passed')
