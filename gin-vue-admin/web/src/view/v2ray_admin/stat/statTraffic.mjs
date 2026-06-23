const UNIT_BYTES = {
  B: 1,
  KB: 1024,
  MB: 1024 * 1024,
  GB: 1024 * 1024 * 1024,
}

const DAY_MS = 24 * 60 * 60 * 1000
const DATE_FORMATTER = new Intl.DateTimeFormat('zh-CN')

export function formatFlow(value) {
  const bytes = Number(value)

  if (!Number.isFinite(bytes) || bytes <= 0) {
    return '0 B'
  }

  if (bytes >= UNIT_BYTES.GB) {
    return `${(bytes / UNIT_BYTES.GB).toFixed(1)} GB`
  }

  if (bytes >= UNIT_BYTES.MB) {
    return `${(bytes / UNIT_BYTES.MB).toFixed(1)} MB`
  }

  if (bytes >= UNIT_BYTES.KB) {
    return `${(bytes / UNIT_BYTES.KB).toFixed(1)} KB`
  }

  return `${bytes.toFixed(1)} B`
}

export function parseTrafficBytes(value) {
  if (typeof value !== 'string') {
    return null
  }

  const match = value.trim().match(/^([\d.]+)\s*(B|KB|MB|GB)$/i)
  if (!match) {
    return null
  }

  const amount = Number.parseFloat(match[1])
  const unit = match[2].toUpperCase()

  if (!Number.isFinite(amount)) {
    return null
  }

  return amount * UNIT_BYTES[unit]
}

export function getTrafficTagType(value) {
  const bytes = parseTrafficBytes(value)

  if (bytes == null) {
    return 'info'
  }

  if (bytes >= UNIT_BYTES.GB) {
    return 'danger'
  }

  if (bytes >= 100 * UNIT_BYTES.MB) {
    return 'warning'
  }

  if (bytes >= 10 * UNIT_BYTES.MB) {
    return 'success'
  }

  return 'info'
}

export function normalizeDateOnlyToUtcIso(value) {
  if (!value) {
    return value
  }

  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }

  return new Date(Date.UTC(date.getFullYear(), date.getMonth(), date.getDate())).toISOString()
}

export function createDefaultTrafficSearchRange({ now = new Date(), days = 30 } = {}) {
  const endDate = now instanceof Date ? new Date(now.getTime()) : new Date(now)
  const rangeDays = Number.isFinite(Number(days)) ? Number(days) : 30
  const startDate = new Date(endDate.getTime() - rangeDays * DAY_MS)

  return {
    startCreatedAt: startDate.toISOString(),
    endCreatedAt: endDate.toISOString(),
  }
}

export function normalizeTrafficSearchRange(searchInfo = {}) {
  const source = searchInfo && typeof searchInfo === 'object' ? searchInfo : {}

  return {
    ...source,
    startCreatedAt: normalizeDateOnlyToUtcIso(source.startCreatedAt),
    endCreatedAt: normalizeDateOnlyToUtcIso(source.endCreatedAt),
  }
}

export function getDateRangeText(searchInfo = {}) {
  const { startCreatedAt, endCreatedAt } = searchInfo

  if (!startCreatedAt || !endCreatedAt) {
    return '最近 30 天'
  }

  const start = DATE_FORMATTER.format(new Date(startCreatedAt))
  const end = DATE_FORMATTER.format(new Date(endCreatedAt))

  return start === end ? start : `${start} - ${end}`
}
