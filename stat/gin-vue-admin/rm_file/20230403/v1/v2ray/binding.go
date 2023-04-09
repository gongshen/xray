package v2ray

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/v2ray"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
    v2rayReq "github.com/flipped-aurora/gin-vue-admin/server/model/v2ray/request"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/service"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type BindingApi struct {
}

var bindingService = service.ServiceGroupApp.V2rayServiceGroup.BindingService


// CreateBinding 创建Binding
// @Tags Binding
// @Summary 创建Binding
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body v2ray.Binding true "创建Binding"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /binding/createBinding [post]
func (bindingApi *BindingApi) CreateBinding(c *gin.Context) {
	var binding v2ray.Binding
	err := c.ShouldBindJSON(&binding)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := bindingService.CreateBinding(&binding); err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteBinding 删除Binding
// @Tags Binding
// @Summary 删除Binding
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body v2ray.Binding true "删除Binding"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /binding/deleteBinding [delete]
func (bindingApi *BindingApi) DeleteBinding(c *gin.Context) {
	var binding v2ray.Binding
	err := c.ShouldBindJSON(&binding)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := bindingService.DeleteBinding(binding); err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteBindingByIds 批量删除Binding
// @Tags Binding
// @Summary 批量删除Binding
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Binding"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /binding/deleteBindingByIds [delete]
func (bindingApi *BindingApi) DeleteBindingByIds(c *gin.Context) {
	var IDS request.IdsReq
    err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := bindingService.DeleteBindingByIds(IDS); err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateBinding 更新Binding
// @Tags Binding
// @Summary 更新Binding
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body v2ray.Binding true "更新Binding"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /binding/updateBinding [put]
func (bindingApi *BindingApi) UpdateBinding(c *gin.Context) {
	var binding v2ray.Binding
	err := c.ShouldBindJSON(&binding)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := bindingService.UpdateBinding(binding); err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindBinding 用id查询Binding
// @Tags Binding
// @Summary 用id查询Binding
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query v2ray.Binding true "用id查询Binding"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /binding/findBinding [get]
func (bindingApi *BindingApi) FindBinding(c *gin.Context) {
	var binding v2ray.Binding
	err := c.ShouldBindQuery(&binding)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if rebinding, err := bindingService.GetBinding(binding.ID); err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rebinding": rebinding}, c)
	}
}

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
