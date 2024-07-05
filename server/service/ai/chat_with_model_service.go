package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
)

var ChatModelServiceApp = new(ChatModelService)

type ChatModelService struct {
}

func (*ChatModelService) Recreation(context *ArticleContext) (*ArticleContext, error) {
	chatModels := global.GVA_CONFIG.ChatModels

	for _, chatModel := range chatModels {
		// 可以通过 WithModel 指定模型
		result, err := ChatServiceApp.Recreation(context, chatModel)
		if err == nil && len(context.Params) > 0 && len(context.Content) > 300 {
			return result, nil
		}
	}

	return nil, nil
}

func (*ChatModelService) HotSpotWrite(account *ai.OfficialAccount, headLine string) (*ArticleContext, error) {

	chatModels := global.GVA_CONFIG.ChatModels

	for _, chatModel := range chatModels {
		// 可以通过 WithModel 指定模型
		result, err := ChatServiceApp.HotSpotWrite(account, headLine, chatModel)
		if err == nil && len(result.Content) > 1000 {
			//
			result.Title = ArticleContentHandlerApp.HandleTitle(result.Title)
			//
			result.Content = ArticleContentHandlerApp.Handle(account, result.Content)

			return result, nil
		}
	}
	return nil, nil
}

func (*ChatModelService) GetKeyWord(title string) string {
	chatModels := global.GVA_CONFIG.ChatModels

	for _, chatModel := range chatModels {
		// 可以通过 WithModel 指定模型
		result := ChatServiceApp.GetKeyWord(title, chatModel)
		if result != "" {
			return result
		}
	}

	return ""
}
