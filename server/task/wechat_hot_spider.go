package task

import "github.com/flipped-aurora/gin-vue-admin/server/service"

func SpiderWechatHotArticle() {
	service.ServiceGroupApp.AIServiceGroup.DajialaService.UpdateHotArticle()
}
