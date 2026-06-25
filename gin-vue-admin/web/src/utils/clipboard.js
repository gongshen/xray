const normalizeClipboardText = (value) => {
  return typeof value === 'string' ? value : String(value ?? '')
}

export const copyTextToClipboard = async(value) => {
  const text = normalizeClipboardText(value)
  if (!text) {
    return false
  }

  try {
    if (globalThis.navigator?.clipboard?.writeText) {
      await globalThis.navigator.clipboard.writeText(text)
      return true
    }
  } catch {
    // Fall back to the textarea path below.
  }

  const doc = globalThis.document
  if (!doc?.body) {
    return false
  }

  const textArea = doc.createElement('textarea')
  textArea.value = text
  textArea.setAttribute('readonly', '')
  textArea.style.position = 'fixed'
  textArea.style.top = '-9999px'
  textArea.style.left = '-9999px'

  doc.body.appendChild(textArea)
  textArea.select()
  textArea.setSelectionRange?.(0, text.length)

  try {
    return Boolean(doc.execCommand?.('copy'))
  } catch {
    return false
  } finally {
    doc.body.removeChild(textArea)
  }
}
