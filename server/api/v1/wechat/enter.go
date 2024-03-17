package wechat

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	WeChatApi
}

var (
	wechatService = service.ServiceGroupApp.AIServiceGroup.WechatService
)
