package ai

import "github.com/flipped-aurora/gin-vue-admin/server/model/ai"

type ChatService interface {
	Recreation(article ai.Article) (*ArticleContext, error)
	HotSpotWrite(topic string) (*ArticleContext, error)
	TopicWrite(topic string) (*ArticleContext, error)
}

var QianfanChat = "qianfan"
var Kimi = "kimi"
var Qianwen = "qianwen"
var AllModel = "all"

var ChatModelServiceApp = new(ChatModelService)

type ChatModelService struct {
}

func (*ChatModelService) Recreation(context *ArticleContext, chatModel string) (*ArticleContext, error) {
	// 可以通过 WithModel 指定模型
	switch chatModel {
	case QianfanChat:
		return KimiServiceApp.Recreation(context)
	case Kimi:
		return KimiServiceApp.Recreation(context)
	case Qianwen:
		return QianwenServiceApp.Recreation(context)
	default:
		result, err := KimiServiceApp.Recreation(context)
		if err == nil && len(context.Params) > 0 {
			return result, nil
		}
		return QianwenServiceApp.Recreation(context)
	}
}

func (*ChatModelService) HotSpotWrite(context *ArticleContext, chatModel string) (*ArticleContext, error) {
	// 可以通过 WithModel 指定模型
	switch chatModel {
	case QianfanChat:
		return KimiServiceApp.HotSpotWrite(context)
	case Kimi:
		return KimiServiceApp.HotSpotWrite(context)
	case Qianwen:
		return QianwenServiceApp.HotSpotWrite(context)
	default:
		result, err := KimiServiceApp.HotSpotWrite(context)
		if err == nil && len(context.Params) > 0 {
			return result, nil
		}
		return QianwenServiceApp.HotSpotWrite(context)
	}
}

func (*ChatModelService) TopicWrite(context *ArticleContext, chatModel string) (*ArticleContext, error) {
	// 可以通过 WithModel 指定模型
	switch chatModel {
	case QianfanChat:
		return KimiServiceApp.TopicWrite(context)
	case Kimi:
		return KimiServiceApp.TopicWrite(context)
	case Qianwen:
		return QianwenServiceApp.TopicWrite(context)
	default:
		result, err := KimiServiceApp.TopicWrite(context)
		if err == nil && len(context.Params) > 0 {
			return result, nil
		}
		return QianwenServiceApp.TopicWrite(context)
	}
}
