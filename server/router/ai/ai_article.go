package ai

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type AIArticleRouter struct {
}

func (e *AIArticleRouter) InitAIArticleRouter(Router *gin.RouterGroup) {
	aiArticleRouter := Router.Group("aiArticle").Use(middleware.OperationRecord())
	aiArticleRouterWithoutRecord := Router.Group("aiArticle")
	aiArticleApi := v1.ApiGroupApp.AIApiGroup.AIArticleApi
	{
		aiArticleRouter.DELETE("aiArticle", aiArticleApi.DeleteAIArticle)                   // 删除文章
		aiArticleRouter.DELETE("deleteAIArticlesByIds", aiArticleApi.DeleteAIArticlesByIds) // 删除文章
		aiArticleRouter.POST("recreation", aiArticleApi.Recreation)                         // 改写文章
	}
	{

		aiArticleRouterWithoutRecord.GET("aiArticle", aiArticleApi.GetAIArticle)         // 获取单一文章信息
		aiArticleRouterWithoutRecord.GET("aiArticleList", aiArticleApi.GetAIArticleList) // 获取文章列表

	}
}
