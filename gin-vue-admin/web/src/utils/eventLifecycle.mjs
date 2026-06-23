export function bindEmitterHandler(emitter, event, handler) {
  emitter.on(event, handler)

  return () => {
    emitter.off(event, handler)
  }
}

export function bindEmitterHandlers(emitter, handlers) {
  const entries = Object.entries(handlers)

  entries.forEach(([event, handler]) => {
    emitter.on(event, handler)
  })

  return () => {
    entries.forEach(([event, handler]) => {
      emitter.off(event, handler)
    })
  }
}

export function bindWindowEvent(target, event, handler) {
  target.addEventListener(event, handler)

  return () => {
    target.removeEventListener(event, handler)
  }
}
