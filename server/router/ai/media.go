package ai

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type MediaRouter struct {
}

func (e *MediaRouter) InitMediaRouter(Router *gin.RouterGroup) {
	mediaRouter := Router.Group("media").Use(middleware.OperationRecord())
	mediaRouterWithoutRecord := Router.Group("media")
	mediaApi := v1.ApiGroupApp.AIApiGroup.MediaApi
	{
		mediaRouter.POST("media", mediaApi.CreateMedia)   // 创建素材
		mediaRouter.DELETE("media", mediaApi.DeleteMedia) // 删除素材
	}
	{
		mediaRouterWithoutRecord.GET("mediaPage", mediaApi.GetMediaPage) // 获取素材列表
	}
}
