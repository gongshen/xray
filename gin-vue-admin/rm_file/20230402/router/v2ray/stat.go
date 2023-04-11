package v2ray

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type StatRouter struct {
}

// InitStatRouter 初始化 Stat 路由信息
func (s *StatRouter) InitStatRouter(Router *gin.RouterGroup) {
	statRouter := Router.Group("stat").Use(middleware.OperationRecord())
	statRouterWithoutRecord := Router.Group("stat")
	var statApi = v1.ApiGroupApp.V2rayApiGroup.StatApi
	{
		statRouter.POST("createStat", statApi.CreateStat)   // 新建Stat
		statRouter.DELETE("deleteStat", statApi.DeleteStat) // 删除Stat
		statRouter.DELETE("deleteStatByIds", statApi.DeleteStatByIds) // 批量删除Stat
		statRouter.PUT("updateStat", statApi.UpdateStat)    // 更新Stat
	}
	{
		statRouterWithoutRecord.GET("findStat", statApi.FindStat)        // 根据ID获取Stat
		statRouterWithoutRecord.GET("getStatList", statApi.GetStatList)  // 获取Stat列表
	}
}
