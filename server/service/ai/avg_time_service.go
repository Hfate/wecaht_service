package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
)

type AvgTimeService struct {
}

var AvgTimeServiceApp = new(AvgTimeService)

func (*AvgTimeService) FindAvgTime() int64 {
	var avgTime int64
	global.GVA_DB.Model(&ai.AvgTime{}).Where("id=1").Select("span_time").Find(&avgTime)
	return avgTime
}

func (*AvgTimeService) UpdateAvgTime(spanTime int64) {
	global.GVA_DB.Model(&ai.AvgTime{}).Where("id=1").Update("span_time", spanTime)
}
