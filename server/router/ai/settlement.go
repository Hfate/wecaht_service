package ai

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type SettlementRouter struct {
}

func (e *SettlementRouter) InitSettlementRouter(Router *gin.RouterGroup) {
	settlementRouterWithoutRecord := Router.Group("settlement")
	settlementApi := v1.ApiGroupApp.AIApiGroup.SettlementApi
	{
		settlementRouterWithoutRecord.GET("list", settlementApi.GetSettlementList) // 获取结算列表
		settlementRouterWithoutRecord.GET("download", settlementApi.Download)      // 获取结算列表

	}
}
