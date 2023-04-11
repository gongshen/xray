package v2ray

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	v2rayReq "github.com/flipped-aurora/gin-vue-admin/server/model/v2ray/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
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
	jwtId := utils.GetUserID(c)
	pageInfo.Category = "user"
	fmt.Println("1212")
	pageInfo.Tag = strconv.Itoa(int(jwtId))
	fmt.Println("3333")
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
	jwtId := utils.GetUserID(c)
	pageInfo.Category = "user"
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
