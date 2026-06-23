function numberOrFallback(value, fallback) {
  const number = Number(value)
  return Number.isFinite(number) ? number : fallback
}

export function normalizeStatTableResponse(response, fallback = {}) {
  const fallbackPage = numberOrFallback(fallback.page, 1)
  const fallbackPageSize = numberOrFallback(fallback.pageSize, 10)

  if (response?.code !== 0) {
    return {
      ok: false,
      list: [],
      total: 0,
      page: fallbackPage,
      pageSize: fallbackPageSize,
      message: response?.msg || '',
    }
  }

  const data = response.data || {}
  const list = Array.isArray(data.list) ? [...data.list] : []

  return {
    ok: true,
    list,
    total: numberOrFallback(data.total, 0),
    page: numberOrFallback(data.page, fallbackPage),
    pageSize: numberOrFallback(data.pageSize, fallbackPageSize),
    message: '',
  }
}
