package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/silenceper/wechat/v2/officialaccount/server"

	aiModel "github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	wechatApi "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/draft"
	"github.com/silenceper/wechat/v2/officialaccount/material"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"go.uber.org/zap"
	"net/http"
)

type WechatService struct {
}

var WechatServiceApp = &WechatService{}
var memoryCache = cache.NewMemory()

var wc = wechatApi.NewWechat()

func init() {
	wc.SetCache(memoryCache)
}

func (*WechatService) ServeWechat(rw http.ResponseWriter, req *http.Request) {

	serverList := WechatServiceApp.ServeList(rw, req)

	// 找到一个可以折腾的公众号
	for _, s := range serverList {
		//处理消息接收以及回复
		err := s.Serve()
		if err != nil {
			global.GVA_LOG.Error("ServeWechat", zap.String("err", err.Error()))
			continue
		}
		//发送回复的消息
		err = s.Send()
		if err != nil {
			global.GVA_LOG.Error("ServeWechat", zap.String("err", err.Error()))
			return
		}
	}

}

func (*WechatService) ServeList(rw http.ResponseWriter, req *http.Request) []*server.Server {
	list, _ := OfficialAccountServiceApp.List()

	result := make([]*server.Server, 0)
	for _, item := range list {
		cfg := &offConfig.Config{
			AppID:          item.AppId,
			AppSecret:      item.AppSecret,
			Token:          item.Token,
			EncodingAESKey: item.EncodingAESKey,
		}
		officialAccount := wc.GetOfficialAccount(cfg)

		// 传入request和responseWriter
		s := officialAccount.GetServer(req, rw)
		//设置接收消息的处理方法
		s.SetMessageHandler(reply)

		result = append(result, s)
	}

	return result
}

func (*WechatService) AddMaterial(dbOfficialAccount aiModel.OfficialAccount, fileName string) (mediaID string, url string, err error) {
	cfg := &offConfig.Config{
		AppID:          dbOfficialAccount.AppId,
		AppSecret:      dbOfficialAccount.AppSecret,
		Token:          dbOfficialAccount.Token,
		EncodingAESKey: dbOfficialAccount.EncodingAESKey,
	}
	officialAccount := wc.GetOfficialAccount(cfg)
	// 获取素材API
	m := officialAccount.GetMaterial()
	mediaID, url, err = m.AddMaterial(material.MediaTypeImage, fileName)

	return
}

func (*WechatService) PublishArticle(dbOfficialAccount aiModel.OfficialAccount, aiArticle aiModel.AIArticle) (publishId int64,
	mediaID string, msgId, msgDataID int64, err error) {

	cfg := &offConfig.Config{
		AppID:          dbOfficialAccount.AppId,
		AppSecret:      dbOfficialAccount.AppSecret,
		Token:          dbOfficialAccount.Token,
		EncodingAESKey: dbOfficialAccount.EncodingAESKey,
	}
	officialAccount := wc.GetOfficialAccount(cfg)

	imageMedia := MediaServiceApp.RandomByAccountId(dbOfficialAccount.AppId)

	// 获取草稿箱api
	d := officialAccount.GetDraft()
	mediaID, err = d.AddDraft([]*draft.Article{{
		Title:        aiArticle.Title,
		ThumbMediaID: imageMedia.MediaID,
		//ThumbURL:     "https://mmbiz.qpic.cn/sz_mmbiz_jpg/uO29ibicRxJ0QfibQSCptBtjsyia61jSn4V7RRX8aLcMUwN7adJhfyaj788qibHVibnOicDyeTAWAor7GGDP6fz1N499A/640?wx_fmt=webp&amp",
		Author: dbOfficialAccount.DefaultAuthorName,
		//Digest:       "test",
		ShowCoverPic: 1,
		Content:      aiArticle.Content,
	}})
	if err != nil {
		return 0, "", 0, 0, err
	}
	//global.GVA_LOG.Info("PublishArticle AddDraft:", zap.String("mediaID", mediaID))
	//
	//p := officialAccount.GetBroadcast()
	//result, err := p.SendNews(nil, mediaID, false)
	//global.GVA_LOG.Info("PublishArticle SendNews:", zap.String("result", utils.Parse2Json(result)))
	//
	//if err != nil {
	//	// 群发不行  试试单发
	//	freePublish := officialAccount.GetFreePublish()
	//	publishID, err := freePublish.Publish(mediaID)
	//	global.GVA_LOG.Info("PublishArticle Publish:", zap.Int64("publishID", publishID), zap.Error(err))
	//	return publishID, mediaID, result.MsgID, result.MsgDataID, err
	//}

	return 0, mediaID, 0, 0, err
}

func reply(msg *message.MixMessage) *message.Reply {
	//回复消息：演示回复用户发送的消息
	text := message.NewText(msg.Content)
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
}
