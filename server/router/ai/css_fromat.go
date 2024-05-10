package ai

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CssFormatRouter struct {
}

func (e *CssFormatRouter) InitCssFormatRouter(Router *gin.RouterGroup) {
	cssFormatRouter := Router.Group("cssFormat").Use(middleware.OperationRecord())
	cssFormatRouterWithoutRecord := Router.Group("cssFormat")
	cssFormatApi := v1.ApiGroupApp.AIApiGroup.CssFormatApi
	{
		cssFormatRouter.POST("cssFormat", cssFormatApi.CreateCssFormat)   // 创建排版
		cssFormatRouter.PUT("cssFormat", cssFormatApi.UpdateCssFormat)    // 更新排版
		cssFormatRouter.DELETE("cssFormat", cssFormatApi.DeleteCssFormat) // 删除排版

	}
	{
		cssFormatRouterWithoutRecord.GET("cssFormat", cssFormatApi.GetCssFormat)         // 获取单一排版信息
		cssFormatRouterWithoutRecord.GET("cssFormatList", cssFormatApi.GetCssFormatList) // 获取排版列表
	}
}
