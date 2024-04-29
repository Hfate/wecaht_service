package utils

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/panjf2000/ants"
	"go.uber.org/zap"
)

var bizPool *ants.Pool

func init() {
	// 初始化协程池   每次仅允许1个并发上传下载
	bizPool, _ = ants.NewPool(4, ants.WithPanicHandler(panicHandler))
}

func BizPool() *ants.Pool {
	return bizPool
}

func panicHandler(err interface{}) {
	global.GVA_LOG.Error("pool panic", zap.Any("err", err))
}
