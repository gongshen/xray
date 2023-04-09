import service from '@/utils/request'

// @Tags Binding
// @Summary 创建Binding
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Binding true "创建Binding"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /binding/createBinding [post]
export const createBinding = (data) => {
  return service({
    url: '/binding/createBinding',
    method: 'post',
    data
  })
}

// @Tags Binding
// @Summary 删除Binding
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Binding true "删除Binding"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /binding/deleteBinding [delete]
export const deleteBinding = (data) => {
  return service({
    url: '/binding/deleteBinding',
    method: 'delete',
    data
  })
}

// @Tags Binding
// @Summary 删除Binding
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Binding"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /binding/deleteBinding [delete]
export const deleteBindingByIds = (data) => {
  return service({
    url: '/binding/deleteBindingByIds',
    method: 'delete',
    data
  })
}

// @Tags Binding
// @Summary 更新Binding
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Binding true "更新Binding"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /binding/updateBinding [put]
export const updateBinding = (data) => {
  return service({
    url: '/binding/updateBinding',
    method: 'put',
    data
  })
}

// @Tags Binding
// @Summary 用id查询Binding
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.Binding true "用id查询Binding"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /binding/findBinding [get]
export const findBinding = (params) => {
  return service({
    url: '/binding/findBinding',
    method: 'get',
    params
  })
}

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
    url: '/binding/getBindingList',
    method: 'get',
    params
  })
}

export const shareBinding = (params) => {
  return service({
    url: '/binding/shareBinding',
    method: 'get',
    params
  })
}