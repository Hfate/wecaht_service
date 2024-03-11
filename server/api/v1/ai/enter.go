package ai

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	PortalApi
	ArticleApi
	BenchmarkAccountApi
	OfficialAccountApi
	HotspotApi
	AIArticleApi
}

var (
	portalService           = service.ServiceGroupApp.AIServiceGroup.PortalService
	aiArticleService        = service.ServiceGroupApp.AIServiceGroup.AIArticleService
	articleService          = service.ServiceGroupApp.AIServiceGroup.ArticleService
	benchmarkAccountService = service.ServiceGroupApp.AIServiceGroup.BenchmarkAccountService
	officialAccountService  = service.ServiceGroupApp.AIServiceGroup.OfficialAccountService
	wxTokenService          = service.ServiceGroupApp.AIServiceGroup.WxTokenService
	hotspotService          = service.ServiceGroupApp.AIServiceGroup.HotspotService
)
