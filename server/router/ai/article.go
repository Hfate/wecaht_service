package ai

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ArticleRouter struct {
}

func (e *ArticleRouter) InitArticleRouter(Router *gin.RouterGroup) {
	articleRouter := Router.Group("article").Use(middleware.OperationRecord())
	articleRouterWithoutRecord := Router.Group("article")
	articleApi := v1.ApiGroupApp.AIApiGroup.ArticleApi
	{
		articleRouter.DELETE("article", articleApi.DeleteArticle)                   // 删除文章
		articleRouter.DELETE("deleteArticlesByIds", articleApi.DeleteArticlesByIds) // 删除文章
		articleRouter.POST("recreation", articleApi.Recreation)                     // 改写文章
		articleRouter.POST("upload", articleApi.UploadArticle)                      // 改写文章
	}
	{
		articleRouterWithoutRecord.GET("template", articleApi.Template)          // 获取单一文章信息
		articleRouterWithoutRecord.GET("article", articleApi.GetArticle)         // 获取单一文章信息
		articleRouterWithoutRecord.GET("articleList", articleApi.GetArticleList) // 获取文章列表
		articleRouterWithoutRecord.GET("download", articleApi.Download)          // 下载文章
	}
}
