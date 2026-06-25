let disposeDomMouseDown = null

export const initDom = () => {
  if (import.meta.env.MODE !== 'development' || disposeDomMouseDown) {
    return
  }

  const handleMouseDown = (e) => {
    if (e.shiftKey && e.button === 0) {
      e.preventDefault()
      sendRequestToOpenFileInEditor(getFilePath(e))
    }
  }

  document.addEventListener('mousedown', handleMouseDown)
  disposeDomMouseDown = () => {
    document.removeEventListener('mousedown', handleMouseDown)
    disposeDomMouseDown = null
  }
}

export const disposeDom = () => {
  disposeDomMouseDown?.()
}

const getFilePath = (e) => {
  let element = e
  if (e.target) {
    element = e.target
  }
  if (!element || !element.getAttribute) return null
  if (element.getAttribute('code-location')) {
    return element.getAttribute('code-location')
  }
  return getFilePath(element.parentNode)
}

const sendRequestToOpenFileInEditor = (filePath) => {
  if (!filePath) {
    return
  }

  const protocol = window.location.protocol
    ? window.location.protocol
    : 'http:'
  const hostname = window.location.hostname
    ? window.location.hostname
    : 'localhost'
  const port = window.location.port ? window.location.port : '80'
  fetch(`${protocol}//${hostname}:${port}/gvaPositionCode?filePath=${encodeURIComponent(filePath)}`)
    .catch(() => {})
}
