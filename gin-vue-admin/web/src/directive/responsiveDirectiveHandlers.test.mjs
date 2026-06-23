import assert from 'node:assert/strict'
import {
  bindElementResizeHandler,
  unbindElementResizeHandler,
} from './responsiveDirectiveHandlers.mjs'

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
  const target = createTarget()
  const element = {}
  const handler = () => {}

  bindElementResizeHandler(element, 'responsive-table', handler, target)
  unbindElementResizeHandler(element, 'responsive-table', target)

  assert.equal(target.calls.length, 2)
  assert.equal(target.calls[0].type, 'add')
  assert.equal(target.calls[0].event, 'resize')
  assert.equal(target.calls[0].handler, handler)
  assert.equal(target.calls[1].type, 'remove')
  assert.equal(target.calls[1].event, 'resize')
  assert.equal(target.calls[1].handler, handler)
}

{
  const target = createTarget()
  const element = {}
  const firstHandler = () => {}
  const secondHandler = () => {}

  bindElementResizeHandler(element, 'responsive-form', firstHandler, target)
  bindElementResizeHandler(element, 'responsive-form', secondHandler, target)

  assert.equal(target.calls.length, 3)
  assert.equal(target.calls[0].type, 'add')
  assert.equal(target.calls[0].handler, firstHandler)
  assert.equal(target.calls[1].type, 'remove')
  assert.equal(target.calls[1].handler, firstHandler)
  assert.equal(target.calls[2].type, 'add')
  assert.equal(target.calls[2].handler, secondHandler)
}

{
  const target = createTarget()
  unbindElementResizeHandler({}, 'missing', target)

  assert.equal(target.calls.length, 0)
}

console.log('responsiveDirectiveHandlers tests passed')
