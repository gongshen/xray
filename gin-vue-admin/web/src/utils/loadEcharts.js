let echartsPromise = null

export const loadEcharts = () => {
  if (!echartsPromise) {
    echartsPromise = import('echarts').catch((error) => {
      echartsPromise = null
      throw error
    })
  }

  return echartsPromise
}