const normalizePort = (value) => {
  const port = Number(value)
  return Number.isFinite(port) ? port : 0
}

const formatPort = (value) => {
  const port = normalizePort(value)
  return port > 0 ? String(port) : '默认端口'
}

export const getServerPortChangeWarnings = (original, current) => {
  if (!original || !current) {
    return []
  }

  const warnings = []
  if (normalizePort(original.port) !== normalizePort(current.port)) {
    warnings.push(`端口已从 ${formatPort(original.port)} 改为 ${formatPort(current.port)}，请在相应服务器上同步修改 Xray 程序配置中的端口。`)
  }

  if (normalizePort(original.stat_port) !== normalizePort(current.stat_port)) {
    warnings.push(`统计端口已从 ${formatPort(original.stat_port)} 改为 ${formatPort(current.stat_port)}，请在相应服务器上同步修改 stat 的端口，不然流量就无法上报。`)
  }

  return warnings
}

export const buildServerPortChangeReminder = (original, current) => {
  const warnings = getServerPortChangeWarnings(original, current)
  if (!warnings.length) {
    return ''
  }

  return `检测到端口配置变更：\n${warnings.map((warning, index) => `${index + 1}. ${warning}`).join('\n')}`
}
