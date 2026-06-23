import assert from 'node:assert/strict'
import {
  bindLayoutEventHandlers,
  getLayoutState,
} from './layoutEvents.mjs'

function createEmitter() {
  const calls = []
  return {
    calls,
    on(event, handler) {
      calls.push({ type: 'on', event, handler })
    },
    off(event, handler) {
      calls.push({ type: 'off', event, handler })
    },
  }
}

function createTarget() {
  const calls = []
  return {
    calls,
    addEventListener(event, handler) {
      calls.push({ type: 'add', event, handler })
    },
    removeEventListener(event, handler) {
      calls.push({ type: 'remove', event, handler })
    },
  }
}

assert.deepEqual(getLayoutState(999), {
  isMobile: true,
  isSider: false,
  isCollapse: true,
})
assert.deepEqual(getLayoutState(1000), {
  isMobile: false,
  isSider: false,
  isCollapse: true,
})
assert.deepEqual(getLayoutState(1199), {
  isMobile: false,
  isSider: false,
  isCollapse: true,
})
assert.deepEqual(getLayoutState(1200), {
  isMobile: false,
  isSider: true,
  isCollapse: false,
})

const emitter = createEmitter()
const target = createTarget()
const handlers = {
  onReload: () => {},
  onShowLoading: () => {},
  onCloseLoading: () => {},
  onResize: () => {},
}

const dispose = bindLayoutEventHandlers({
  emitter,
  target,
  ...handlers,
})

dispose()

assert.deepEqual(emitter.calls.map(call => [call.type, call.event, call.handler]), [
  ['on', 'reload', handlers.onReload],
  ['on', 'showLoading', handlers.onShowLoading],
  ['on', 'closeLoading', handlers.onCloseLoading],
  ['off', 'reload', handlers.onReload],
  ['off', 'showLoading', handlers.onShowLoading],
  ['off', 'closeLoading', handlers.onCloseLoading],
])
assert.deepEqual(target.calls.map(call => [call.type, call.event, call.handler]), [
  ['add', 'resize', handlers.onResize],
  ['remove', 'resize', handlers.onResize],
])

console.log('layoutEvents tests passed')
