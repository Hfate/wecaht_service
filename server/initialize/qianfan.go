package initialize

import (
	"github.com/baidubce/bce-qianfan-sdk/go/qianfan"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

func QianFan() {
	qianFanCfg := global.GVA_CONFIG.QianFan
	qianfan.GetConfig().AccessKey = qianFanCfg.AccessKey
	qianfan.GetConfig().SecretKey = qianFanCfg.SecretKey
}
