package v2ray

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type StatRouter struct {
}

// InitStatRouter 初始化 Stat 路由信息
func (s *StatRouter) InitStatRouter(Router *gin.RouterGroup) {
	statRouter := Router.Group("v2ray/stat").Use(middleware.OperationRecord())
	var statApi = v1.ApiGroupApp.V2rayApiGroup.StatApi
	{
		statRouter.GET("getStatList", statApi.GetStatList)     // 获取Stat列表
		statRouter.GET("getStatCharts", statApi.GetStatCharts) // 获取流量统计图表展示
	}
}
