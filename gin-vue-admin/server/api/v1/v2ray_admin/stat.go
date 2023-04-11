package v2ray_admin

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/v2ray"
	v2rayReq "github.com/flipped-aurora/gin-vue-admin/server/model/v2ray/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type StatApi struct {
}

var statService = service.ServiceGroupApp.V2rayAdminServiceGroup.StatService

// CreateStat 创建Stat
// @Tags Stat
// @Summary 创建Stat
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body v2ray_admin.Stat true "创建Stat"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /stat/createStat [post]
func (statApi *StatApi) CreateStat(c *gin.Context) {
	var stat v2ray.Stat
	err := c.ShouldBindJSON(&stat)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	verify := utils.Rules{
		"Category": {utils.NotEmpty()},
	}
	if err := utils.Verify(stat, verify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := statService.CreateStat(&stat); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteStat 删除Stat
// @Tags Stat
// @Summary 删除Stat
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body v2ray_admin.Stat true "删除Stat"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /stat/deleteStat [delete]
func (statApi *StatApi) DeleteStat(c *gin.Context) {
	var stat v2ray.Stat
	err := c.ShouldBindJSON(&stat)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := statService.DeleteStat(stat); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteStatByIds 批量删除Stat
// @Tags Stat
// @Summary 批量删除Stat
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Stat"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /stat/deleteStatByIds [delete]
func (statApi *StatApi) DeleteStatByIds(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := statService.DeleteStatByIds(IDS); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateStat 更新Stat
// @Tags Stat
// @Summary 更新Stat
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body v2ray_admin.Stat true "更新Stat"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /stat/updateStat [put]
func (statApi *StatApi) UpdateStat(c *gin.Context) {
	var stat v2ray.Stat
	err := c.ShouldBindJSON(&stat)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	verify := utils.Rules{
		"Category": {utils.NotEmpty()},
	}
	if err := utils.Verify(stat, verify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := statService.UpdateStat(stat); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindStat 用id查询Stat
// @Tags Stat
// @Summary 用id查询Stat
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query v2ray.Stat true "用id查询Stat"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /stat/findStat [get]
func (statApi *StatApi) FindStat(c *gin.Context) {
	var stat v2ray.Stat
	err := c.ShouldBindQuery(&stat)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if restat, err := statService.GetStat(stat.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"restat": restat}, c)
	}
}

// GetStatList 分页获取Stat列表
// @Tags Stat
// @Summary 分页获取Stat列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query v2rayReq.StatSearch true "分页获取Stat列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /stat/getStatList [get]
func (statApi *StatApi) GetStatList(c *gin.Context) {
	var pageInfo v2rayReq.StatSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := statService.GetStatInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

func (statApi *StatApi) GetStatCharts(c *gin.Context) {
	var pageInfo v2rayReq.StatSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, err := statService.GetStatCharts(&pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(list, "获取成功", c)
	}
}
