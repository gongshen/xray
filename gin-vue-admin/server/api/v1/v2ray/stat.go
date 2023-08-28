package v2ray

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	v2rayReq "github.com/flipped-aurora/gin-vue-admin/server/model/v2ray/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type StatApi struct {
}

var statService = service.ServiceGroupApp.V2rayAdminServiceGroup.StatService

func (statApi *StatApi) GetStatCharts(c *gin.Context) {
	var pageInfo v2rayReq.StatSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if pageInfo.EndCreatedAt == nil || pageInfo.StartCreatedAt == nil {
		response.FailWithMessage("请输入时间", c)
		return
	}
	if pageInfo.EndCreatedAt.Sub(*pageInfo.StartCreatedAt) > time.Hour*24*365 {
		response.FailWithMessage("查询时间不能超过一年", c)
		return
	}
	jwtId := utils.GetUserID(c)
	pageInfo.Tag = strconv.Itoa(int(jwtId))
	if list, err := statService.GetStatCharts(&pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(list, "获取成功", c)
	}
}

func (statApi *StatApi) GetStatList(c *gin.Context) {
	var pageInfo v2rayReq.StatSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if pageInfo.EndCreatedAt == nil || pageInfo.StartCreatedAt == nil {
		response.FailWithMessage("请输入时间", c)
		return
	}
	if pageInfo.EndCreatedAt.Sub(*pageInfo.StartCreatedAt) > time.Hour*24*365 {
		response.FailWithMessage("查询时间不能超过一年", c)
		return
	}
	jwtId := utils.GetUserID(c)
	pageInfo.Tag = strconv.Itoa(int(jwtId))
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
