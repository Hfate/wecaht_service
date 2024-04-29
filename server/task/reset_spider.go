package task

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"gorm.io/gorm"
)

func ResetSpider(db *gorm.DB) {
	db.Model(&ai.BenchmarkAccount{}).Where("1 = 1").Update("spider_flag", 0)
}
