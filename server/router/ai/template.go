package ai

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type TemplateRouter struct {
}

func (e *TemplateRouter) InitTemplateRouter(Router *gin.RouterGroup) {
	templateRouter := Router.Group("template").Use(middleware.OperationRecord())
	templateRouterWithoutRecord := Router.Group("template")
	templateApi := v1.ApiGroupApp.AIApiGroup.TemplateApi
	{
		templateRouter.DELETE("template", templateApi.DeleteTemplate) // 删除模板
		templateRouter.POST("update", templateApi.UpdateTemplate)     // 更新模板
		templateRouter.POST("clone", templateApi.CloneTemplate)       // 克隆模板
	}
	{
		templateRouterWithoutRecord.GET("template", templateApi.GetTemplate)         // 获取单一模板信息
		templateRouterWithoutRecord.GET("templateList", templateApi.GetTemplateList) // 获取模板列表

	}
}
