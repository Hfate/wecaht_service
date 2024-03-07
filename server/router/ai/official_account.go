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
		officialAccountRouter.POST("officialAccount", officialAccountApi.CreateOfficialAccount)   // 创建门户
		officialAccountRouter.PUT("officialAccount", officialAccountApi.UpdateOfficialAccount)    // 更新门户
		officialAccountRouter.DELETE("officialAccount", officialAccountApi.DeleteOfficialAccount) // 删除门户
	}
	{
		officialAccountRouterWithoutRecord.GET("officialAccount", officialAccountApi.GetOfficialAccount)         // 获取单一门户信息
		officialAccountRouterWithoutRecord.GET("officialAccountList", officialAccountApi.GetOfficialAccountList) // 获取门户列表
	}
}
