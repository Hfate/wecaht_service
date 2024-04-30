package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/spf13/cast"
)

type AvgTimeService struct {
}

var AvgTimeServiceApp = new(AvgTimeService)

func (*AvgTimeService) FindAvgTime() int64 {
	var avgTime float64
	global.GVA_DB.Model(&ai.AvgTime{}).Where("1=1").Select("avg(span_time)").Find(&avgTime)
	return cast.ToInt64(avgTime)
}

func (*AvgTimeService) UpdateAvgTime(spanTime int64) {
	newAvgTime := ai.AvgTime{
		ID:       utils.GenID(),
		SpanTime: spanTime,
	}
	global.GVA_DB.Model(&ai.AvgTime{}).Create(newAvgTime)
}
