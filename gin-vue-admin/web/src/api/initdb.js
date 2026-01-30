import service from '@/utils/request'

export const initDB = (data) => {
  return service({
    url: '/init/initdb',
    method: 'post',
    data
  })
}

export const checkDB = () => {
  return service({
    url: '/init/checkdb',
    method: 'post'
  })
}
