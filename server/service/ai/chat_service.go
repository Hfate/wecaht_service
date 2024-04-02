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
		result, err := KimiServiceApp.Recreation(article)
		if err == nil {
			return result, nil
		}
		return QianfanServiceApp.Recreation(article)
	}

	return nil, nil
}

func (*ChatModelService) HotSpotWrite(link string, chatModel string) (*ArticleContext, error) {
	// 可以通过 WithModel 指定模型
	switch chatModel {
	case QianfanChat:
		return QianfanServiceApp.HotSpotWrite(link)
	case Kimi:

	default:
		result, err := KimiServiceApp.HotSpotWrite(link)
		if err == nil {
			return result, nil
		}
		return QianfanServiceApp.HotSpotWrite(link)
	}

	return nil, nil
}

func (*ChatModelService) TopicWrite(topic, chatModel string) (*ArticleContext, error) {
	// 可以通过 WithModel 指定模型
	switch chatModel {
	case QianfanChat:
		return QianfanServiceApp.TopicWrite(topic)
	case Kimi:

	default:
		result, err := KimiServiceApp.TopicWrite(topic)
		if err == nil {
			return result, nil
		}
		return QianfanServiceApp.TopicWrite(topic)
	}
	return nil, nil
}
