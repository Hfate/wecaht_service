package ai

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type HotspotRouter struct {
}

func (e *HotspotRouter) InitHotspotRouter(Router *gin.RouterGroup) {
	hotspotRouter := Router.Group("hotspot").Use(middleware.OperationRecord())
	hotspotRouterWithoutRecord := Router.Group("hotspot")
	hotspotApi := v1.ApiGroupApp.AIApiGroup.HotspotApi
	{

		hotspotRouter.DELETE("hotspot", hotspotApi.DeleteHotspot) // 删除门户
	}
	{

		hotspotRouterWithoutRecord.GET("hotspotList", hotspotApi.GetHotspotList) // 获取门户列表
	}
}
