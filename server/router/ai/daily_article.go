package ai

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type DailyArticleRouter struct {
}

func (e *DailyArticleRouter) InitDailyArticleRouter(Router *gin.RouterGroup) {
	dailyArticleRouter := Router.Group("dailyArticle").Use(middleware.OperationRecord())
	dailyArticleRouterWithoutRecord := Router.Group("dailyArticle")
	dailyArticleApi := v1.ApiGroupApp.AIApiGroup.DailyArticleApi
	{
		dailyArticleRouter.DELETE("dailyArticle", dailyArticleApi.DeleteDailyArticle)              // 删除文章
		dailyArticleRouter.DELETE("deleteArticlesByIds", dailyArticleApi.DeleteDailyArticlesByIds) // 删除文章
		dailyArticleRouter.POST("recreation", dailyArticleApi.Recreation)                          // 改写文章
	}
	{
		dailyArticleRouterWithoutRecord.GET("dailyArticle", dailyArticleApi.GetDailyArticle)         // 获取单一文章信息
		dailyArticleRouterWithoutRecord.GET("dailyArticleList", dailyArticleApi.GetDailyArticleList) // 获取文章列表
	}
}
