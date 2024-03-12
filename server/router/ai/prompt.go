package ai

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type PromptRouter struct {
}

func (e *PromptRouter) InitPromptRouter(Router *gin.RouterGroup) {
	promptRouter := Router.Group("prompt").Use(middleware.OperationRecord())
	promptRouterWithoutRecord := Router.Group("prompt")
	promptApi := v1.ApiGroupApp.AIApiGroup.PromptApi
	{
		promptRouter.POST("prompt", promptApi.CreatePrompt)   // 创建Prompt
		promptRouter.PUT("prompt", promptApi.UpdatePrompt)    // 更新Prompt
		promptRouter.DELETE("prompt", promptApi.DeletePrompt) // 删除Prompt
	}
	{
		promptRouterWithoutRecord.GET("prompt", promptApi.GetPrompt)         // 获取单一Prompt信息
		promptRouterWithoutRecord.GET("promptList", promptApi.GetPromptList) // 获取Prompt列表
	}
}
