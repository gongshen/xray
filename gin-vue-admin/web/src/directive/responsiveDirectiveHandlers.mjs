const handlersByElement = new WeakMap()

function getHandlers(element) {
  let handlers = handlersByElement.get(element)

  if (!handlers) {
    handlers = new Map()
    handlersByElement.set(element, handlers)
  }

  return handlers
}

export function bindElementResizeHandler(element, key, handler, target = window) {
  unbindElementResizeHandler(element, key, target)

  const handlers = getHandlers(element)
  handlers.set(key, handler)
  target.addEventListener('resize', handler)
}

export function unbindElementResizeHandler(element, key, target = window) {
  const handlers = handlersByElement.get(element)

  if (!handlers?.has(key)) {
    return
  }

  const handler = handlers.get(key)
  target.removeEventListener('resize', handler)
  handlers.delete(key)

  if (handlers.size === 0) {
    handlersByElement.delete(element)
  }
}
