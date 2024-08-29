package ai

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	aiModel "github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	wechatApi "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/draft"
	"github.com/silenceper/wechat/v2/officialaccount/freepublish"
	"github.com/silenceper/wechat/v2/officialaccount/material"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/silenceper/wechat/v2/officialaccount/server"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type WechatService struct {
}

var WechatServiceApp = &WechatService{}
var memoryCache = cache.NewMemory()

var wc = wechatApi.NewWechat()

func init() {
	wc.SetCache(memoryCache)
}

func (*WechatService) PublisherSettlement() {
	list, _ := OfficialAccountServiceApp.List()

	global.GVA_LOG.Info("PublisherSettlement", zap.Int("account-length", len(list)))

	wechatSettlementList := make([]*aiModel.WechatSettlement, 0)
	for _, item := range list {
		cfg := &offConfig.Config{
			AppID:          item.AppId,
			AppSecret:      item.AppSecret,
			Token:          item.Token,
			EncodingAESKey: item.EncodingAESKey,
		}
		officialAccount := wc.GetOfficialAccount(cfg)
		settlementList, err := officialAccount.GetDataCube().GetPublisherSettlement("2024-01-01", "2024-09-01", 1, 100)
		if err != nil {
			global.GVA_LOG.Error("PublisherSettlement", zap.String("AccountName", item.AccountName), zap.Any("err", err))
			continue
		}

		global.GVA_LOG.Info("PublisherSettlement", zap.String("AccountName", item.AccountName),
			zap.String("resp", utils.Parse2Json(settlementList)), zap.Int("length-settlementList", len(settlementList.SettlementList)))

		if len(settlementList.SettlementList) > 0 {
			for _, set := range settlementList.SettlementList {
				wechatSettlementList = append(wechatSettlementList, &aiModel.WechatSettlement{
					BASEMODEL:      BaseModel(),
					AccountName:    item.AccountName,
					AccountId:      item.AccountId,
					Date:           set.Date,
					Zone:           set.Zone,
					Month:          set.Month,
					Order:          set.Order,
					SettStatus:     set.SettStatus,
					SettledRevenue: set.SettledRevenue,
					SettNo:         set.SettNo,
					MailSendCnt:    set.MailSendCnt,
					SlotRevenue:    utils.Parse2Json(set.SlotRevenue),
				})
			}

			global.GVA_DB.Where("1=1").Where("account_id=?", item.AccountId).Delete(&aiModel.WechatSettlement{})

			err2 := global.GVA_DB.Model(&aiModel.WechatSettlement{}).Create(wechatSettlementList).Error
			if err2 != nil {
				global.GVA_LOG.Error("PublisherSettlement", zap.String("err", err2.Error()))
				continue
			}
		}

	}
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

func (*WechatService) AddMaterial(dbOfficialAccount *aiModel.OfficialAccount, fileName string) (mediaID string, url string, err error) {
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

func (*WechatService) ImageUpload(dbOfficialAccount *aiModel.OfficialAccount, fileName string) (url string, err error) {
	cfg := &offConfig.Config{
		AppID:          dbOfficialAccount.AppId,
		AppSecret:      dbOfficialAccount.AppSecret,
		Token:          dbOfficialAccount.Token,
		EncodingAESKey: dbOfficialAccount.EncodingAESKey,
	}
	officialAccount := wc.GetOfficialAccount(cfg)
	// 获取素材API
	m := officialAccount.GetMaterial()
	url, err = m.ImageUpload(fileName)

	return
}

func (ws *WechatService) PublishArticle(dbOfficialAccount *aiModel.OfficialAccount, aiArticleList []aiModel.AIArticle) (publishId int64,
	mediaID string, msgId, msgDataID int64, err error) {

	cfg := &offConfig.Config{
		AppID:          dbOfficialAccount.AppId,
		AppSecret:      dbOfficialAccount.AppSecret,
		Token:          dbOfficialAccount.Token,
		EncodingAESKey: dbOfficialAccount.EncodingAESKey,
	}
	officialAccount := wc.GetOfficialAccount(cfg)

	draftList := make([]*draft.Article, 0)
	for _, aiArticle := range aiArticleList {
		if aiArticle.ArticleStatus == 1 {
			continue
		}

		// 搜索封面托片
		filePath := utils.SearchAndSave(aiArticle.Title)

		if filePath == "" {
			// 找不到 则使用ai
			keyword := ChatModelServiceApp.GetKeyWord(aiArticleList[0].Title)
			filePath = utils.SearchAndSave(keyword)
			if filePath == "" {
				global.GVA_LOG.Error("没有找到封面图片", zap.String("title", aiArticle.Title),
					zap.String("appId", dbOfficialAccount.AppId),
					zap.String("appName", dbOfficialAccount.AccountName))
			}
		}

		imgMediaID, _, err2 := MediaServiceApp.CreateMediaByPath(dbOfficialAccount.AppId, filePath)
		if err2 != nil {
			// 默认图片
			media := MediaServiceApp.FindByAccountId(dbOfficialAccount.AppId, 1)
			if media != nil {
				imgMediaID = media.MediaID
			} else {
				global.GVA_LOG.Error("上传封面图片失败", zap.String("title", aiArticle.Title),
					zap.String("filePath", filePath),
					zap.String("appId", dbOfficialAccount.AppId),
					zap.String("appName", dbOfficialAccount.AccountName), zap.Error(err2))
			}
		}

		// 替换富文本图片
		content := ws.replaceImg(dbOfficialAccount.AppId, aiArticle.Content)

		draftList = append(draftList, &draft.Article{
			Title:        aiArticle.Title,
			ThumbMediaID: imgMediaID,
			//ThumbURL:     "https://mmbiz.qpic.cn/sz_mmbiz_jpg/uO29ibicRxJ0QfibQSCptBtjsyia61jSn4V7RRX8aLcMUwN7adJhfyaj788qibHVibnOicDyeTAWAor7GGDP6fz1N499A/640?wx_fmt=webp&amp",
			Author: dbOfficialAccount.DefaultAuthorName,
			//Digest:       "test",
			NeedOpenComment: dbOfficialAccount.NeedOpenComment,
			ShowCoverPic:    1,
			Content:         content,
		})

	}

	if len(draftList) == 0 {
		return 0, "", 0, 0, errors.New("全部发送草稿失败")
	}

	if len(draftList) > 5 {
		draftList = draftList[0:5]
	}

	// 获取草稿箱api
	d := officialAccount.GetDraft()
	mediaID, err = d.AddDraft(draftList)

	fmt.Println("appId："+dbOfficialAccount.AppId+";appName:"+dbOfficialAccount.AccountName+";发布草稿：mediaID:", mediaID, err)
	if err != nil {
		global.GVA_LOG.Error("PublishArticle AddDraft:", zap.String("err", err.Error()),
			zap.String("appId", officialAccount.GetContext().AppID),
			zap.String("dbAppId", dbOfficialAccount.AppId),
			zap.String("mediaID", mediaID),
		)

		return 0, "", 0, 0, err
	}
	global.GVA_LOG.Info("PublishArticle AddDraft:", zap.String("mediaID", mediaID))

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

func (ws *WechatService) replaceImg(appId, content string) string {
	// 使用 goquery.NewDocumentFromReader 从字符串创建一个 Document 实例
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		global.GVA_LOG.Error("replaceImg", zap.Error(err))
	}

	// 使用 Find 查找所有的 img 标签
	doc.Find("img").Each(func(index int, img *goquery.Selection) {
		// 获取现有的 src 属性值
		originSrc, _ := img.Attr("src")

		// 上传至公众号
		newSrc, err2 := MediaServiceApp.ImageUpload(appId, originSrc)

		global.GVA_LOG.Info("替换公众号配图",
			zap.String("originSrc", originSrc),
			zap.String("URL", newSrc),
			zap.String("appId", appId), zap.Error(err2))

		// 设置新的 src 属性值
		img.SetAttr("src", newSrc)

	})

	htmlResult, _ := doc.Html()

	// 输出修改后的 HTML 文档
	return htmlResult
}

func reply(msg *message.MixMessage) *message.Reply {
	//回复消息：演示回复用户发送的消息
	text := message.NewText(msg.Content)
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
}

func (*WechatService) BatchGetHistoryArticleList(dbOfficialAccount *aiModel.OfficialAccount) (list freepublish.ArticleList, err error) {
	cfg := &offConfig.Config{
		AppID:          dbOfficialAccount.AppId,
		AppSecret:      dbOfficialAccount.AppSecret,
		Token:          dbOfficialAccount.Token,
		EncodingAESKey: dbOfficialAccount.EncodingAESKey,
	}
	officialAccount := wc.GetOfficialAccount(cfg)
	// 获取发布接口
	p := officialAccount.GetFreePublish()
	// 获取最近五篇文章
	return p.Paginate(0, 5, true)
}
