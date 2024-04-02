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
		hotspotRouter.POST("create", hotspotApi.CreateArticle)    // 删除头条
		hotspotRouter.DELETE("hotspot", hotspotApi.DeleteHotspot) // 删除头条
	}
	{

		hotspotRouterWithoutRecord.GET("hotspotList", hotspotApi.GetHotspotList) // 获取头条列表
	}
}
