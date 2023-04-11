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
	bindingRouterWithoutRecord := Router.Group("v2ray/binding").Use(middleware.OperationRecord())
	var bindingApi = v1.ApiGroupApp.V2rayApiGroup.BindingApi
	{
		bindingRouterWithoutRecord.GET("getBindingList", bindingApi.GetBindingList) // 获取Binding列表
		bindingRouterWithoutRecord.GET("shareBinding", bindingApi.ShareBinding)     // 分享binding
	}
}
