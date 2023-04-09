import service from '@/utils/request'

// @Tags Server
// @Summary 创建Server
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Server true "创建Server"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /server/createServer [post]
export const createServer = (data) => {
  return service({
    url: '/server/createServer',
    method: 'post',
    data
  })
}

// @Tags Server
// @Summary 删除Server
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Server true "删除Server"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /server/deleteServer [delete]
export const deleteServer = (data) => {
  return service({
    url: '/server/deleteServer',
    method: 'delete',
    data
  })
}

// @Tags Server
// @Summary 删除Server
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Server"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /server/deleteServer [delete]
export const deleteServerByIds = (data) => {
  return service({
    url: '/server/deleteServerByIds',
    method: 'delete',
    data
  })
}

// @Tags Server
// @Summary 更新Server
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Server true "更新Server"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /server/updateServer [put]
export const updateServer = (data) => {
  return service({
    url: '/server/updateServer',
    method: 'put',
    data
  })
}

// @Tags Server
// @Summary 用id查询Server
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.Server true "用id查询Server"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /server/findServer [get]
export const findServer = (params) => {
  return service({
    url: '/server/findServer',
    method: 'get',
    params
  })
}

// @Tags Server
// @Summary 分页获取Server列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取Server列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /server/getServerList [get]
export const getServerList = (params) => {
  return service({
    url: '/server/getServerList',
    method: 'get',
    params
  })
}

export const getAllServerApi = () => {
  return service({
    url: '/server/getAllServer',
    method: 'post'
  })
}