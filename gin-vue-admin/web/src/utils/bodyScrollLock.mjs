const bodyScrollLocks = new WeakMap()

export function lockBodyScroll(target = globalThis.document?.body) {
  if (!target?.style) {
    return () => {}
  }

  const current = bodyScrollLocks.get(target) || {
    count: 0,
    overflow: target.style.overflow,
  }

  current.count += 1
  bodyScrollLocks.set(target, current)
  target.style.overflow = 'hidden'

  let disposed = false
  return () => {
    if (disposed) {
      return
    }

    disposed = true
    const latest = bodyScrollLocks.get(target)
    if (!latest) {
      return
    }

    latest.count -= 1
    if (latest.count <= 0) {
      target.style.overflow = latest.overflow
      bodyScrollLocks.delete(target)
    } else {
      bodyScrollLocks.set(target, latest)
    }
  }
}
