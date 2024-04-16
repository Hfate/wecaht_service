package task

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"go.uber.org/zap"
)

func GenerateArticle() {
	err := service.ServiceGroupApp.AIServiceGroup.AIArticleService.GenerateDailyArticle()
	if err != nil {
		global.GVA_LOG.Error("生成文章失败", zap.Error(err))
	}
}
