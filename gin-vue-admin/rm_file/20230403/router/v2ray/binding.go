package v2ray

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type BindingRouter struct {
}

// InitBindingRouter 初始化 Binding 路由信息
func (s *BindingRouter) InitBindingRouter(Router *gin.RouterGroup) {
	bindingRouter := Router.Group("binding").Use(middleware.OperationRecord())
	bindingRouterWithoutRecord := Router.Group("binding")
	var bindingApi = v1.ApiGroupApp.V2rayApiGroup.BindingApi
	{
		bindingRouter.POST("createBinding", bindingApi.CreateBinding)   // 新建Binding
		bindingRouter.DELETE("deleteBinding", bindingApi.DeleteBinding) // 删除Binding
		bindingRouter.DELETE("deleteBindingByIds", bindingApi.DeleteBindingByIds) // 批量删除Binding
		bindingRouter.PUT("updateBinding", bindingApi.UpdateBinding)    // 更新Binding
	}
	{
		bindingRouterWithoutRecord.GET("findBinding", bindingApi.FindBinding)        // 根据ID获取Binding
		bindingRouterWithoutRecord.GET("getBindingList", bindingApi.GetBindingList)  // 获取Binding列表
	}
}
