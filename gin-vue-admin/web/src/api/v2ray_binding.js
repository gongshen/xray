import service from '@/utils/request'

export const getBindingList = (params) => {
  return service({
    url: '/v2ray/binding/getBindingList',
    method: 'get',
    params
  })
}

export const shareBinding = (params) => {
  return service({
    url: '/v2ray/binding/shareBinding',
    method: 'get',
    params
  })
}
