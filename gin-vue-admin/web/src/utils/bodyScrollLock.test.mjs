import assert from 'node:assert/strict'
import { lockBodyScroll } from './bodyScrollLock.mjs'

const target = { style: { overflow: 'auto' } }

const releaseFirst = lockBodyScroll(target)
assert.equal(target.style.overflow, 'hidden')

const releaseSecond = lockBodyScroll(target)
assert.equal(target.style.overflow, 'hidden')

releaseFirst()
assert.equal(target.style.overflow, 'hidden')

releaseSecond()
assert.equal(target.style.overflow, 'auto')

releaseSecond()
assert.equal(target.style.overflow, 'auto')

assert.doesNotThrow(() => lockBodyScroll(null)())
assert.doesNotThrow(() => lockBodyScroll({})())

console.log('bodyScrollLock tests passed')
