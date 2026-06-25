const HISTORY_STORAGE_KEY = 'historys'
const ACTIVE_VALUE_STORAGE_KEY = 'activeValue'
const NEED_CLOSE_ALL_STORAGE_KEY = 'needCloseAll'

const getStorage = () => globalThis.sessionStorage || null

export function readHistoryStorage(fallback = []) {
  try {
    const raw = getStorage()?.getItem(HISTORY_STORAGE_KEY)
    if (!raw) {
      return fallback
    }
    const parsed = JSON.parse(raw)
    return Array.isArray(parsed) ? parsed : fallback
  } catch {
    return fallback
  }
}

export function writeHistoryStorage(history) {
  try {
    getStorage()?.setItem(HISTORY_STORAGE_KEY, JSON.stringify(history || []))
  } catch {
    // Session storage can fail in private mode or when quota is exceeded.
  }
}

export function readActiveValue() {
  try {
    return getStorage()?.getItem(ACTIVE_VALUE_STORAGE_KEY) || ''
  } catch {
    return ''
  }
}

export function writeActiveValue(value) {
  try {
    getStorage()?.setItem(ACTIVE_VALUE_STORAGE_KEY, value || '')
  } catch {
    // Keep in-memory navigation usable even if storage is unavailable.
  }
}

export function shouldCloseAllHistory() {
  try {
    return getStorage()?.getItem(NEED_CLOSE_ALL_STORAGE_KEY) === 'true'
  } catch {
    return false
  }
}

export function markCloseAllHistory() {
  try {
    getStorage()?.setItem(NEED_CLOSE_ALL_STORAGE_KEY, 'true')
  } catch {
    // Keep authority switching usable even if session storage is unavailable.
  }
}

export function clearCloseAllHistoryFlag() {
  try {
    getStorage()?.removeItem(NEED_CLOSE_ALL_STORAGE_KEY)
  } catch {
    // Ignore storage failures.
  }
}
