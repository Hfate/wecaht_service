package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"go.uber.org/zap"
)

type PoolService struct {
}

var PoolServiceApp = &PoolService{}

func (fb *PoolService) SubmitBizTask(task func()) {
	err := utils.BizPool().Submit(task)
	if err != nil {
		global.GVA_LOG.Error("SubmitBizTask", zap.Error(err))
	}
}
