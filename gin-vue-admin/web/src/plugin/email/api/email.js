import service from '@/utils/request'

export const emailTest = (data) => {
  return service({
    url: '/email/emailTest',
    method: 'post',
    data
  })
}

export const sendEmail = (data) => {
  return service({
    url: '/email/sendEmail',
    method: 'post',
    data
  })
}
