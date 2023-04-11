import service from '@/utils/request'

// @Tags Binding
// @Summary 分页获取Binding列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取Binding列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /binding/getBindingList [get]
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