export function getLayoutState(screenWidth) {
  if (screenWidth < 1000) {
    return {
      isMobile: true,
      isSider: false,
      isCollapse: true,
    }
  }

  if (screenWidth < 1200) {
    return {
      isMobile: false,
      isSider: false,
      isCollapse: true,
    }
  }

  return {
    isMobile: false,
    isSider: true,
    isCollapse: false,
  }
}

export function bindLayoutEventHandlers({
  emitter,
  target = window,
  onReload,
  onShowLoading,
  onCloseLoading,
  onResize,
}) {
  const emitterHandlers = [
    ['reload', onReload],
    ['showLoading', onShowLoading],
    ['closeLoading', onCloseLoading],
  ]

  emitterHandlers.forEach(([event, handler]) => {
    emitter.on(event, handler)
  })
  target.addEventListener('resize', onResize)

  return () => {
    emitterHandlers.forEach(([event, handler]) => {
      emitter.off(event, handler)
    })
    target.removeEventListener('resize', onResize)
  }
}
