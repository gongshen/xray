import service from '@/utils/request'

// @Tags Stat
// @Summary 创建Stat
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Stat true "创建Stat"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /stat/createStat [post]
export const createStat = (data) => {
  return service({
    url: '/v2ray_admin/stat/createStat',
    method: 'post',
    data
  })
}

// @Tags Stat
// @Summary 删除Stat
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Stat true "删除Stat"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /stat/deleteStat [delete]
export const deleteStat = (data) => {
  return service({
    url: '/v2ray_admin/stat/deleteStat',
    method: 'delete',
    data
  })
}

// @Tags Stat
// @Summary 删除Stat
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Stat"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /stat/deleteStat [delete]
export const deleteStatByIds = (data) => {
  return service({
    url: '/v2ray_admin/stat/deleteStatByIds',
    method: 'delete',
    data
  })
}

// @Tags Stat
// @Summary 更新Stat
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Stat true "更新Stat"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /stat/updateStat [put]
export const updateStat = (data) => {
  return service({
    url: '/v2ray_admin/stat/updateStat',
    method: 'put',
    data
  })
}

// @Tags Stat
// @Summary 用id查询Stat
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.Stat true "用id查询Stat"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /stat/findStat [get]
export const findStat = (params) => {
  return service({
    url: '/v2ray_admin/stat/findStat',
    method: 'get',
    params
  })
}

// @Tags Stat
// @Summary 分页获取Stat列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取Stat列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /stat/getStatList [get]
export const getStatList = (params) => {
  return service({
    url: '/v2ray_admin/stat/getStatList',
    method: 'get',
    params
  })
}

// @Tags Stat
// @Summary 获取Stat图表数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取Stat列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /stat/getStatCharts [get]
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