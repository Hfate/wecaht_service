package ai

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type PortalRouter struct {
}

func (e *PortalRouter) InitPortalRouter(Router *gin.RouterGroup) {
	portalRouter := Router.Group("portal").Use(middleware.OperationRecord())
	portalRouterWithoutRecord := Router.Group("portal")
	portalApi := v1.ApiGroupApp.AIApiGroup.PortalApi
	{
		portalRouter.POST("portal", portalApi.CreatePortal)   // 创建门户
		portalRouter.PUT("portal", portalApi.UpdatePortal)    // 更新门户
		portalRouter.DELETE("portal", portalApi.DeletePortal) // 删除门户
	}
	{
		portalRouterWithoutRecord.GET("portal", portalApi.GetPortal)         // 获取单一门户信息
		portalRouterWithoutRecord.GET("portalList", portalApi.GetPortalList) // 获取门户列表
	}
}
