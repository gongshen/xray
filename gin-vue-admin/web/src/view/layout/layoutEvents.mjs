import {
  bindEmitterHandlers,
  bindWindowEvent,
} from '../../utils/eventLifecycle.mjs'

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
  const disposeEmitterHandlers = bindEmitterHandlers(emitter, {
    reload: onReload,
    showLoading: onShowLoading,
    closeLoading: onCloseLoading,
  })
  const disposeResize = bindWindowEvent(target, 'resize', onResize)

  return () => {
    disposeEmitterHandlers()
    disposeResize()
  }
}
