const HTTP_PROTOCOLS = new Set(['http:', 'https:'])

export const normalizeExternalUrl = (value) => {
  return typeof value === 'string' ? value.trim() : ''
}

export const isExternalHttpUrl = (value) => {
  const url = normalizeExternalUrl(value)
  if (!url) {
    return false
  }

  try {
    return HTTP_PROTOCOLS.has(new URL(url).protocol)
  } catch {
    return false
  }
}

export const openExternalUrl = (value) => {
  const url = normalizeExternalUrl(value)
  if (!isExternalHttpUrl(url)) {
    return false
  }

  window.open(url, '_blank', 'noopener,noreferrer')
  return true
}
