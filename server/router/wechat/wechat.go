package wechat

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type WeChatRouter struct {
}

func (e *WeChatRouter) InitWeChatRouter(Router *gin.RouterGroup) {

	wechatRouterWithoutRecord := Router.Group("wechat")
	wechatApi := v1.ApiGroupApp.WechatGroup.WeChatApi
	{
		wechatRouterWithoutRecord.GET("callback", wechatApi.CallBack) // 微信回调
	}
}
