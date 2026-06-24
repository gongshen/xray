import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync('src/view/system/state.vue', 'utf8')
const mountedStart = source.indexOf('onMounted(() => {')
const unmountedStart = source.indexOf('onUnmounted(() => {')

assert.notEqual(mountedStart, -1, 'state page should register mounted lifecycle work')
assert.notEqual(unmountedStart, -1, 'state page should register unmounted lifecycle cleanup')

const mountedBlock = source.slice(mountedStart, unmountedStart)
const topLevelBetweenHooks = source.slice(source.indexOf('})', mountedStart) + 2, unmountedStart)

assert.match(mountedBlock, /timer\.value\s*=\s*setInterval\(/, 'refresh timer should start during mounted lifecycle')
assert.doesNotMatch(topLevelBetweenHooks, /timer\.value\s*=\s*setInterval\(/, 'refresh timer should not start at script setup top level')
assert.match(source.slice(unmountedStart), /clearInterval\(timer\.value\)/, 'refresh timer should be cleared during unmount')

console.log('stateTimerLifecycle tests passed')
