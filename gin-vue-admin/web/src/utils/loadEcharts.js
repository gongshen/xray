let echartsPromise = null

export const loadEcharts = () => {
  if (!echartsPromise) {
    echartsPromise = import('./echartsCore.js').catch((error) => {
      echartsPromise = null
      throw error
    })
  }

  return echartsPromise
}
