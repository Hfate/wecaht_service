package ai

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type TopicRouter struct {
}

func (e *TopicRouter) InitTopicRouter(Router *gin.RouterGroup) {
	topicRouter := Router.Group("topic").Use(middleware.OperationRecord())
	topicRouterWithoutRecord := Router.Group("topic")
	topicApi := v1.ApiGroupApp.AIApiGroup.TopicApi
	{
		topicRouter.POST("topic", topicApi.CreateTopic)   // 创建主题
		topicRouter.PUT("topic", topicApi.UpdateTopic)    // 更新主题
		topicRouter.DELETE("topic", topicApi.DeleteTopic) // 删除主题
	}
	{
		topicRouterWithoutRecord.GET("topicList", topicApi.GetTopicList) // 获取主题列表
		topicRouterWithoutRecord.GET("topicPage", topicApi.GetTopicPage) // 获取主题列表
	}
}
