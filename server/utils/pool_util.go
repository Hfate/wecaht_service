package utils

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/panjf2000/ants"
	"go.uber.org/zap"
)

var bizPool *ants.Pool

func init() {
	//
	bizPool, _ = ants.NewPool(2, ants.WithPanicHandler(panicHandler))
}

func BizPool() *ants.Pool {
	return bizPool
}

func panicHandler(err interface{}) {
	global.GVA_LOG.Error("pool panic", zap.Any("err", err))
}
