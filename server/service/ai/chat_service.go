package ai

import "github.com/flipped-aurora/gin-vue-admin/server/model/ai"

type ChatService interface {
	Recreation(article ai.Article) (*ArticleContext, error)
	HotSpotWrite(topic string) (*ArticleContext, error)
	TopicWrite(topic string) (*ArticleContext, error)
}

var QianfanChat = "qianfan"
var Kimi = "kimi"
var AllModel = "all"

var ChatModelServiceApp = new(ChatModelService)

type ChatModelService struct {
}

func (*ChatModelService) Recreation(article ai.Article, chatModel string) (*ArticleContext, error) {
	// 可以通过 WithModel 指定模型
	switch chatModel {
	case QianfanChat:
		return QianfanServiceApp.Recreation(article)
	case Kimi:

	default:
		return QianfanServiceApp.Recreation(article)
	}

	return nil, nil
}

func (*ChatModelService) HotSpotWrite(topic string, chatModel string) (*ArticleContext, error) {
	// 可以通过 WithModel 指定模型
	switch chatModel {
	case QianfanChat:
		return QianfanServiceApp.HotSpotWrite(topic)
	case Kimi:

	default:
		return QianfanServiceApp.HotSpotWrite(topic)
	}

	return nil, nil
}

func (*ChatModelService) TopicWrite(topic, chatModel string) (*ArticleContext, error) {
	// 可以通过 WithModel 指定模型
	switch chatModel {
	case QianfanChat:
		return QianfanServiceApp.TopicSpotWrite(topic)
	case Kimi:

	default:
		return QianfanServiceApp.TopicSpotWrite(topic)
	}
	return nil, nil
}
