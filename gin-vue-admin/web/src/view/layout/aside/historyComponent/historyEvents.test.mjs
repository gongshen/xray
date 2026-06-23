import assert from 'node:assert/strict'
import {
  bindBodyClickHandler,
  bindEmitterHandlers,
} from './historyEvents.mjs'

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

{
  const target = createTarget()
  const handler = () => {}
  const dispose = bindBodyClickHandler(handler, target)

  dispose()

  assert.equal(target.calls.length, 2)
  assert.equal(target.calls[0].type, 'add')
  assert.equal(target.calls[0].event, 'click')
  assert.equal(target.calls[0].handler, handler)
  assert.equal(target.calls[1].type, 'remove')
  assert.equal(target.calls[1].event, 'click')
  assert.equal(target.calls[1].handler, handler)
}

{
  const emitter = createEmitter()
  const handlers = {
    closeThisPage: () => {},
    closeAllPage: () => {},
    mobile: () => {},
    collapse: () => {},
  }
  const dispose = bindEmitterHandlers(emitter, handlers)

  dispose()

  assert.equal(emitter.calls.length, 8)
  Object.entries(handlers).forEach(([event, handler], index) => {
    assert.equal(emitter.calls[index].type, 'on')
    assert.equal(emitter.calls[index].event, event)
    assert.equal(emitter.calls[index].handler, handler)
    assert.equal(emitter.calls[index + 4].type, 'off')
    assert.equal(emitter.calls[index + 4].event, event)
    assert.equal(emitter.calls[index + 4].handler, handler)
  })
}

console.log('historyEvents tests passed')
