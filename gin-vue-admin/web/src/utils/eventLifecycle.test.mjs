import assert from 'node:assert/strict'
import {
  bindEmitterHandler,
  bindEmitterHandlers,
  bindWindowEvent,
} from './eventLifecycle.mjs'

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

{
  const emitter = createEmitter()
  const handler = () => {}
  const dispose = bindEmitterHandler(emitter, 'routeChange', handler)

  dispose()

  assert.deepEqual(emitter.calls.map(call => [call.type, call.event, call.handler]), [
    ['on', 'routeChange', handler],
    ['off', 'routeChange', handler],
  ])
}

{
  const emitter = createEmitter()
  const handlers = {
    mobile: () => {},
    collapse: () => {},
  }
  const dispose = bindEmitterHandlers(emitter, handlers)

  dispose()

  assert.deepEqual(emitter.calls.map(call => [call.type, call.event, call.handler]), [
    ['on', 'mobile', handlers.mobile],
    ['on', 'collapse', handlers.collapse],
    ['off', 'mobile', handlers.mobile],
    ['off', 'collapse', handlers.collapse],
  ])
}

{
  const target = createTarget()
  const handler = () => {}
  const dispose = bindWindowEvent(target, 'resize', handler)

  dispose()

  assert.deepEqual(target.calls.map(call => [call.type, call.event, call.handler]), [
    ['add', 'resize', handler],
    ['remove', 'resize', handler],
  ])
}

console.log('eventLifecycle tests passed')
