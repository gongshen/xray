import service from '@/utils/request'

export const createServer = (data) => {
  return service({
    url: '/v2ray_admin/server/createServer',
    method: 'post',
    data
  })
}

export const deleteServer = (data) => {
  return service({
    url: '/v2ray_admin/server/deleteServer',
    method: 'delete',
    data
  })
}

export const deleteServerByIds = (data) => {
  return service({
    url: '/v2ray_admin/server/deleteServerByIds',
    method: 'delete',
    data
  })
}

export const updateServer = (data) => {
  return service({
    url: '/v2ray_admin/server/updateServer',
    method: 'put',
    data
  })
}

export const findServer = (params) => {
  return service({
    url: '/v2ray_admin/server/findServer',
    method: 'get',
    params
  })
}

export const getServerList = (params) => {
  return service({
    url: '/v2ray_admin/server/getServerList',
    method: 'get',
    params
  })
}

export const getAllServerApi = () => {
  return service({
    url: '/v2ray_admin/server/getAllServer',
    method: 'post'
  })
}

export const restartXrayApi = (data) => {
  return service({
    url: '/v2ray_admin/server/restartXray',
    method: 'put',
    data
  })
}

export const restartVPSApi = (data) => {
  return service({
    url: '/v2ray_admin/server/restartVPS',
    method: 'post',
    data
  })
}

export const analyzeUserTrafficApi = (data) => {
  return service({
    url: '/v2ray_admin/server/analyzeUserTraffic',
    method: 'post',
    data
  })
}

export const classifyTrafficTargetsApi = (data) => {
  return service({
    url: '/v2ray_admin/server/classifyTrafficTargets',
    method: 'post',
    data
  })
}
