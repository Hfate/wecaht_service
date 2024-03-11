package ai

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type BenchmarkAccountRouter struct {
}

func (e *BenchmarkAccountRouter) InitBenchmarkAccountRouter(Router *gin.RouterGroup) {
	benchmarkAccountRouter := Router.Group("benchmark").Use(middleware.OperationRecord())
	benchmarkAccountRouterWithoutRecord := Router.Group("benchmark")
	wxTokenRouterWithoutRecord := Router.Group("wxToken")
	benchmarkAccountApi := v1.ApiGroupApp.AIApiGroup.BenchmarkAccountApi
	{
		wxTokenRouterWithoutRecord.PUT("wxToken", benchmarkAccountApi.UpdateWxToken)           // 更新wx token
		benchmarkAccountRouter.POST("benchmark", benchmarkAccountApi.CreateBenchmarkAccount)   // 创建对标账号
		benchmarkAccountRouter.DELETE("benchmark", benchmarkAccountApi.DeleteBenchmarkAccount) // 删除对标账号
	}
	{
		benchmarkAccountRouterWithoutRecord.GET("benchmark", benchmarkAccountApi.GetBenchmarkAccount)         // 获取单一对标账号信息
		benchmarkAccountRouterWithoutRecord.GET("benchmarkList", benchmarkAccountApi.GetBenchmarkAccountList) // 获取对标账号列表
	}
}
