const BYTE_UNITS = ['B', 'KB', 'MB', 'GB', 'TB']

function finiteNumber(value) {
  const number = Number(value)
  return Number.isFinite(number) ? number : 0
}

export function calculatePercent(used, total) {
  const usedValue = finiteNumber(used)
  const totalValue = finiteNumber(total)

  if (totalValue <= 0) {
    return 0
  }

  return Math.round((usedValue / totalValue) * 100)
}

export function formatTime(timestamp) {
  const seconds = finiteNumber(timestamp)

  if (seconds <= 0) {
    return '-'
  }

  return new Date(seconds * 1000).toLocaleString('zh-CN')
}

export function formatBytes(value) {
  let bytes = finiteNumber(value)

  if (bytes <= 0) {
    return '0 B'
  }

  let unitIndex = 0
  while (bytes >= 1024 && unitIndex < BYTE_UNITS.length - 1) {
    bytes /= 1024
    unitIndex++
  }

  return `${bytes.toFixed(2)} ${BYTE_UNITS[unitIndex]}`
}

export function formatSize(value) {
  const mb = finiteNumber(value)

  if (mb <= 0) {
    return '0 GB'
  }

  if (mb >= 1024) {
    return `${(mb / 1024).toFixed(1)} GB`
  }

  return `${mb} MB`
}

export function todayCompact(date = new Date()) {
  const current = date instanceof Date ? date : new Date(date)
  const month = `${current.getMonth() + 1}`.padStart(2, '0')
  const day = `${current.getDate()}`.padStart(2, '0')

  return `${current.getFullYear()}${month}${day}`
}

export function parseClockMinute(value) {
  const match = `${value}`.trim().match(/^(\d{1,2}):(\d{2})$/)

  if (!match) {
    return null
  }

  const hour = Number(match[1])
  const minute = Number(match[2])

  if (hour < 0 || hour > 23 || minute < 0 || minute > 59) {
    return null
  }

  return hour * 60 + minute
}

export function validateTrafficAnalysisQuery({ currentServer, form, maxMinutes = 120 } = {}) {
  const userTag = String(form?.user_tag || '').trim()
  const date = String(form?.date || '').trim()
  const start = parseClockMinute(form?.start)
  const end = parseClockMinute(form?.end)

  if (!currentServer?.ID) {
    return '服务器信息无效'
  }

  if (!userTag) {
    return '请选择用户'
  }

  if (!/^\d{8}$/.test(date)) {
    return '日期格式应为 20260617'
  }

  if (start === null || end === null) {
    return '时间格式应为 8:10 或 09:00'
  }

  if (end < start) {
    return '结束时间不能早于开始时间'
  }

  if (end - start > maxMinutes) {
    return '时间范围不能超过2小时'
  }

  return ''
}

export function getTrafficAnalysisTotals(rows = []) {
  return rows.reduce((totals, row) => {
    totals.down += finiteNumber(row?.down)
    totals.up += finiteNumber(row?.up)
    totals.total += finiteNumber(row?.total)
    return totals
  }, { down: 0, up: 0, total: 0 })
}

export function targetNames(targets = []) {
  if (!Array.isArray(targets) || targets.length === 0) {
    return []
  }

  return [...new Set(targets.map(item => item.target).filter(Boolean))]
}

export function rowTrafficAnalysisTargets(row) {
  const targets = new Set()

  targetNames(row?.targets).forEach(target => {
    const value = String(target || '').trim().toLowerCase()
    if (value) {
      targets.add(value)
    }
  })

  return [...targets]
}

export function formatUserOption(user = {}) {
  const nickName = user.nickName || user.userName || '未命名用户'
  return `${nickName} (ID:${user.ID})`
}

export function isServerOnline(timestamp, thresholdSeconds, nowSeconds = Math.floor(Date.now() / 1000)) {
  const updatedAt = finiteNumber(timestamp)

  if (updatedAt <= 0) {
    return false
  }

  return nowSeconds - updatedAt < thresholdSeconds
}
