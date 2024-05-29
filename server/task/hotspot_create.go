package task

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	ai2 "github.com/flipped-aurora/gin-vue-admin/server/service/ai"
	"gorm.io/gorm"
	"time"
)

func HotspotCreate(db *gorm.DB) error {
	// 获取当前时间
	now := time.Now()

	// 设置时间为当天的开始，即00:00:00
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	hotspotList := make([]*ai.Hotspot, 0)

	accountList := make([]*ai.OfficialAccount, 0)
	db.Model(&ai.OfficialAccount{}).Where("topic=?", "热点").Find(&accountList)

	limit := len(accountList)

	db.Model(&ai.Hotspot{}).Where("avg_speed>1000000").Where("use_times=0").Where("created_at>?", startOfDay).Limit(limit).Find(&hotspotList)

	for index, hotspot := range hotspotList {

		hotspot.UseTimes = 1
		db.Save(hotspot)

		account := accountList[index]

		result, err := ai2.ChatModelServiceApp.HotSpotWrite(account, hotspot.Headline)

		articleList := make([]ai.AIArticle, 0)
		if err == nil && len(result.Content) > 1000 {
			articleList = append(articleList, ai.AIArticle{
				Title:   result.Title,
				Content: result.Content,
			})

			ai2.WechatServiceApp.PublishArticle(account, articleList)
		}
	}

	return nil
}