package wechat

import (
	"github.com/gin-gonic/gin"
)

type WeChatApi struct{}

// CallBack
// @Tags      Wechat
// @Summary   创建主题
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      wechat.Wechat            true  "主题用户名, 主题手机号码"
// @Success   200   {object}  response.Response{msg=string}  "创建主题"
// @Router    /wechat/callback [post]
func (e *WeChatApi) CallBack(c *gin.Context) {
	wechatService.ServeWechat(c.Writer, c.Request)
}
