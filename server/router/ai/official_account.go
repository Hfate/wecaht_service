package ai

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type OfficialAccountRouter struct {
}

func (e *OfficialAccountRouter) InitOfficialAccountRouter(Router *gin.RouterGroup) {
	officialAccountRouter := Router.Group("officialAccount").Use(middleware.OperationRecord())
	officialAccountRouterWithoutRecord := Router.Group("officialAccount")
	officialAccountApi := v1.ApiGroupApp.AIApiGroup.OfficialAccountApi
	{
		officialAccountRouter.POST("officialAccount", officialAccountApi.CreateOfficialAccount)   // 创建公众号
		officialAccountRouter.PUT("officialAccount", officialAccountApi.UpdateOfficialAccount)    // 更新公众号
		officialAccountRouter.PUT("updateCreateTypes", officialAccountApi.UpdateCreateTypes)      // 更新公众号的创建方式
		officialAccountRouter.DELETE("officialAccount", officialAccountApi.DeleteOfficialAccount) // 删除公众号
		officialAccountRouter.POST("create", officialAccountApi.CreateArticle)                    // 创建文章

	}
	{
		officialAccountRouterWithoutRecord.GET("officialAccount", officialAccountApi.GetOfficialAccount)         // 获取单一公众号信息
		officialAccountRouterWithoutRecord.GET("officialAccountList", officialAccountApi.GetOfficialAccountList) // 获取公众号列表
	}
}
