import service from '@/utils/request'

export const createStat = (data) => {
  return service({
    url: '/v2ray_admin/stat/createStat',
    method: 'post',
    data
  })
}

export const deleteStat = (data) => {
  return service({
    url: '/v2ray_admin/stat/deleteStat',
    method: 'delete',
    data
  })
}

export const deleteStatByIds = (data) => {
  return service({
    url: '/v2ray_admin/stat/deleteStatByIds',
    method: 'delete',
    data
  })
}

export const updateStat = (data) => {
  return service({
    url: '/v2ray_admin/stat/updateStat',
    method: 'put',
    data
  })
}

export const findStat = (params) => {
  return service({
    url: '/v2ray_admin/stat/findStat',
    method: 'get',
    params
  })
}

export const getStatList = (params) => {
  return service({
    url: '/v2ray_admin/stat/getStatList',
    method: 'get',
    params
  })
}

export const getStatCharts = (params) => {
  return service({
    url: '/v2ray_admin/stat/getStatCharts',
    method: 'get',
    params
  })
}

export const getStatRank = (params) => {
  return service({
    url: '/v2ray_admin/stat/getStatRank',
    method: 'get',
    params
  })
}
