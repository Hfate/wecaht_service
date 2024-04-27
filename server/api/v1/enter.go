package v1

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/system"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/wechat"
)

type ApiGroup struct {
	SystemApiGroup system.ApiGroup
	AIApiGroup     ai.ApiGroup
	WechatGroup    wechat.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
