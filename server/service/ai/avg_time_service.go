package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
)

type AvgTimeService struct {
}

var AvgTimeServiceApp = new(AvgTimeService)

func (*AvgTimeService) FindAvgTime() int64 {
	var avgTime int64
	global.GVA_DB.Model(&ai.AvgTime{}).Where("1=1").Select("avg(span_time)").Find(&avgTime)
	return avgTime
}

func (*AvgTimeService) UpdateAvgTime(spanTime int64) {
	newAvgTime := ai.AvgTime{
		ID:       utils.GenID(),
		SpanTime: spanTime,
	}
	global.GVA_DB.Model(&ai.AvgTime{}).Create(newAvgTime)
}
