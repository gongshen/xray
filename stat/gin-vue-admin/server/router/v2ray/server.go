package v2ray

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ServerRouter struct {
}

// InitServerRouter 初始化 Server 路由信息
func (s *ServerRouter) InitServerRouter(Router *gin.RouterGroup) {
	serverRouter := Router.Group("server").Use(middleware.OperationRecord())
	serverRouterWithoutRecord := Router.Group("server")
	var serverApi = v1.ApiGroupApp.V2rayApiGroup.ServerApi
	{
		serverRouter.POST("createServer", serverApi.CreateServer)             // 新建Server
		serverRouter.DELETE("deleteServer", serverApi.DeleteServer)           // 删除Server
		serverRouter.DELETE("deleteServerByIds", serverApi.DeleteServerByIds) // 批量删除Server
		serverRouter.PUT("updateServer", serverApi.UpdateServer)              // 更新Server
	}
	{
		serverRouterWithoutRecord.GET("findServer", serverApi.FindServer)       // 根据ID获取Server
		serverRouterWithoutRecord.GET("getServerList", serverApi.GetServerList) // 获取Server列表
		serverRouterWithoutRecord.POST("getAllServer", serverApi.GetAllServer)  // 获取所有Server
	}
}
