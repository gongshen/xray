import service from '@/utils/request'

export const createStat = (data) => {
  return service({
    url: '/v2ray/stat/createStat',
    method: 'post',
    data
  })
}

export const deleteStat = (data) => {
  return service({
    url: '/v2ray/stat/deleteStat',
    method: 'delete',
    data
  })
}

export const deleteStatByIds = (data) => {
  return service({
    url: '/v2ray/stat/deleteStatByIds',
    method: 'delete',
    data
  })
}

export const updateStat = (data) => {
  return service({
    url: '/v2ray/stat/updateStat',
    method: 'put',
    data
  })
}

export const findStat = (params) => {
  return service({
    url: '/v2ray/stat/findStat',
    method: 'get',
    params
  })
}

export const getStatList = (params) => {
  return service({
    url: '/v2ray/stat/getStatList',
    method: 'get',
    params
  })
}

export const getStatCharts = (params) => {
  return service({
    url: '/v2ray/stat/getStatCharts',
    method: 'get',
    params
  })
}

export const getStatRank = (params) => {
  return service({
    url: '/v2ray/stat/getStatRank',
    method: 'get',
    params
  })
}
