package ai

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	PortalApi
	ArticleApi
	BenchmarkAccountApi
	OfficialAccountApi
}

var (
	portalService           = service.ServiceGroupApp.AIServiceGroup.PortalService
	articleService          = service.ServiceGroupApp.AIServiceGroup.ArticleService
	benchmarkAccountService = service.ServiceGroupApp.AIServiceGroup.BenchmarkAccountService
	officialAccountService  = service.ServiceGroupApp.AIServiceGroup.OfficialAccountService
)
