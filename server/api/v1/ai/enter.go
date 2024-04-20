package ai

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	PortalApi
	ArticleApi
	DailyArticleApi
	BenchmarkAccountApi
	OfficialAccountApi
	HotspotApi
	AIArticleApi
	PromptApi
	TopicApi
	MediaApi
}

var (
	topicService            = service.ServiceGroupApp.AIServiceGroup.TopicService
	portalService           = service.ServiceGroupApp.AIServiceGroup.PortalService
	promptService           = service.ServiceGroupApp.AIServiceGroup.PromptService
	aiArticleService        = service.ServiceGroupApp.AIServiceGroup.AIArticleService
	articleService          = service.ServiceGroupApp.AIServiceGroup.ArticleService
	dailyArticleService     = service.ServiceGroupApp.AIServiceGroup.DailyArticleService
	benchmarkAccountService = service.ServiceGroupApp.AIServiceGroup.BenchmarkAccountService
	officialAccountService  = service.ServiceGroupApp.AIServiceGroup.OfficialAccountService
	wxTokenService          = service.ServiceGroupApp.AIServiceGroup.WxTokenService
	hotspotService          = service.ServiceGroupApp.AIServiceGroup.HotspotService
	mediaService            = service.ServiceGroupApp.AIServiceGroup.MediaService
)
