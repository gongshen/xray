export const fmtTitle = (title = '', now = {}) => {
  const params = now.params || {}
  const query = now.query || {}

  return String(title).replace(/\$\{(.+?)\}/g, (_, key) => {
    const value = params[key] ?? query[key]
    return value == null ? '' : String(value)
  })
}
