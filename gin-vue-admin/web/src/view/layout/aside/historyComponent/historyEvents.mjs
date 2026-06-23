export function bindBodyClickHandler(handler, target = document.body) {
  target.addEventListener('click', handler)

  return () => {
    target.removeEventListener('click', handler)
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
