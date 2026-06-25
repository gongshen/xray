import service from '@/utils/request'

export const createStatApi = (baseUrl, { includeRank = false } = {}) => {
  const api = {
    createStat(data) {
      return service({
        url: `${baseUrl}/createStat`,
        method: 'post',
        data
      })
    },

    deleteStat(data) {
      return service({
        url: `${baseUrl}/deleteStat`,
        method: 'delete',
        data
      })
    },

    deleteStatByIds(data) {
      return service({
        url: `${baseUrl}/deleteStatByIds`,
        method: 'delete',
        data
      })
    },

    updateStat(data) {
      return service({
        url: `${baseUrl}/updateStat`,
        method: 'put',
        data
      })
    },

    findStat(params) {
      return service({
        url: `${baseUrl}/findStat`,
        method: 'get',
        params
      })
    },

    getStatList(params) {
      return service({
        url: `${baseUrl}/getStatList`,
        method: 'get',
        params
      })
    },

    getStatCharts(params) {
      return service({
        url: `${baseUrl}/getStatCharts`,
        method: 'get',
        params
      })
    }
  }

  if (includeRank) {
    api.getStatRank = (params) => {
      return service({
        url: `${baseUrl}/getStatRank`,
        method: 'get',
        params
      })
    }
  }

  return api
}
