package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/service/example"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/wechat"
)

type ServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
	ExampleServiceGroup example.ServiceGroup
	AIServiceGroup      ai.ServiceGroup
	WechatServiceGroup  wechat.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
