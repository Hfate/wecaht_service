package wechat

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/service/ai"
	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"go.uber.org/zap"
	"net/http"
)

type WechatService struct {
}

func (*WechatService) ServeWechat(rw http.ResponseWriter, req *http.Request) {
	wc := wechat.NewWechat()
	//这里本地内存保存access_token，也可选择redis，memcache或者自定cache
	memory := cache.NewMemory()

	dbOfficialAccount, _ := ai.OfficialAccountServiceApp.GetLastOfficialAccount()

	cfg := &offConfig.Config{
		AppID:     dbOfficialAccount.AppId,
		AppSecret: dbOfficialAccount.AppSecret,
		Token:     dbOfficialAccount.Token,
		//EncodingAESKey: "xxxx",
		Cache: memory,
	}
	officialAccount := wc.GetOfficialAccount(cfg)

	// 传入request和responseWriter
	server := officialAccount.GetServer(req, rw)
	//设置接收消息的处理方法
	server.SetMessageHandler(reply)

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		global.GVA_LOG.Error("ServeWechat", zap.String("err", err.Error()))
		return
	}
	//发送回复的消息
	err = server.Send()

	if err != nil {
		global.GVA_LOG.Error("ServeWechat", zap.String("err", err.Error()))
		return
	}
}

func reply(msg *message.MixMessage) *message.Reply {

	//回复消息：演示回复用户发送的消息
	text := message.NewText(msg.Content)
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
}
