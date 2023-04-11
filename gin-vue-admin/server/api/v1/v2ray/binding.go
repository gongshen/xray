package v2ray

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/v2ray"
	v2rayReq "github.com/flipped-aurora/gin-vue-admin/server/model/v2ray/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

type BindingApi struct {
}

var bindingService = service.ServiceGroupApp.V2rayAdminServiceGroup.BindingService

// GetBindingList 分页获取Binding列表
// @Tags Binding
// @Summary 分页获取Binding列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query v2rayReq.BindingSearch true "分页获取Binding列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /binding/getBindingList [get]
func (bindingApi *BindingApi) GetBindingList(c *gin.Context) {
	var pageInfo v2rayReq.BindingSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	jwtId := utils.GetUserID(c)
	pageInfo.UserID = int(jwtId)
	if pageInfo.UserID <= 0 {
		response.FailWithMessage("用户ID不正确", c)
		return
	}
	if list, total, err := bindingService.GetBindingInfoList(pageInfo); err != nil {
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

const vmessSharePattern = `vmess://tcp:%s-%d@%s:%d?type=http#dino2-%s`

// ShareBinding TODO 暂时只支持分享vmess协议
func (bindingApi *BindingApi) ShareBinding(c *gin.Context) {
	var binding v2ray.Binding
	if err := c.ShouldBindQuery(&binding); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	ans, err := bindingService.GetBinding(binding.ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	share2 := defaultV2rayN
	share2.Add = ans.Server.Ip
	share2.Id = ans.User.UUID.String()
	share2.Port = strconv.FormatInt(ans.Server.Port, 10)
	share2.Ps = fmt.Sprintf("dino2-%s", ans.Server.Ip)
	share2.Aid = strconv.FormatInt(ans.AlterID, 10)
	share2Data, err := json.Marshal(share2)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(ShareResponse{
		Share1: fmt.Sprintf(vmessSharePattern, ans.User.UUID.String(), ans.AlterID, ans.Server.Ip, ans.Server.Port, ans.Server.Ip),
		Share2: fmt.Sprintf("vmess://%s", base64.StdEncoding.EncodeToString(share2Data)),
	}, "查询成功", c)
}

type ShareResponse struct {
	Share1 string `json:"share1"`
	Share2 string `json:"share2"`
}

var defaultV2rayN = v2ray.V2ray{
	V:    "2",
	Type: "http",
	Net:  "tcp",
	Scy:  "auto",
}
