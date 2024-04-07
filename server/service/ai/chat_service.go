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

func (*ChatModelService) Recreation(article ai.Article, chatModel string) (*ArticleContext, error) {
	// 可以通过 WithModel 指定模型
	switch chatModel {
	case QianfanChat:
		return KimiServiceApp.Recreation(article)
	case Kimi:
		return KimiServiceApp.Recreation(article)
	case Qianwen:
		return QianwenServiceApp.Recreation(article)
	default:
		result, err := KimiServiceApp.Recreation(article)
		if err == nil {
			return result, nil
		}
		return QianwenServiceApp.Recreation(article)
	}
}

func (*ChatModelService) HotSpotWrite(link string, chatModel string) (*ArticleContext, error) {
	// 可以通过 WithModel 指定模型
	switch chatModel {
	case QianfanChat:
		return KimiServiceApp.HotSpotWrite(link)
	case Kimi:
		return KimiServiceApp.HotSpotWrite(link)
	case Qianwen:
		return QianwenServiceApp.HotSpotWrite(link)
	default:
		result, err := KimiServiceApp.HotSpotWrite(link)
		if err == nil {
			return result, nil
		}
		return QianwenServiceApp.HotSpotWrite(link)
	}
}

func (*ChatModelService) TopicWrite(topic, chatModel string) (*ArticleContext, error) {
	// 可以通过 WithModel 指定模型
	switch chatModel {
	case QianfanChat:
		return KimiServiceApp.TopicWrite(topic)
	case Kimi:
		return KimiServiceApp.TopicWrite(topic)
	case Qianwen:
		return QianwenServiceApp.TopicWrite(topic)
	default:
		result, err := KimiServiceApp.TopicWrite(topic)
		if err == nil {
			return result, nil
		}
		return QianwenServiceApp.TopicWrite(topic)
	}
}
