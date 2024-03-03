package initialize

import (
	"github.com/songzhibin97/gkit/cache/local_cache"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
)

func OtherInit() {
	dr, err := utils.ParseDuration(global.GVA_CONFIG.JWT.ExpiresTime)
	if err != nil {
		panic(err)
	}
	_, err = utils.ParseDuration(global.GVA_CONFIG.JWT.BufferTime)
	if err != nil {
		panic(err)
	}

	utils.Init("2024-02-29", 1)

	global.BlackCache = local_cache.NewCache(
		local_cache.SetDefaultExpire(dr),
	)

}
