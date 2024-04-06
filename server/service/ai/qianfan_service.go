package ai

//
//import (
//	"bytes"
//	"context"
//	"encoding/json"
//	"github.com/baidubce/bce-qianfan-sdk/go/qianfan"
//	"github.com/flipped-aurora/gin-vue-admin/server/global"
//	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
//	"github.com/flipped-aurora/gin-vue-admin/server/utils"
//	"go.uber.org/zap"
//	"strings"
//	"text/template"
//)
//
//type QianfanService struct {
//}
//
//var QianfanServiceApp = new(QianfanService)
//
//func (*QianfanService) HotSpotWrite(topic string) (*ArticleContext, error) {
//	chat := qianfan.NewChatCompletion(qianfan.WithModel("ERNIE-Bot-4"))
//	articleContext := &ArticleContext{}
//	articleContext.Topic = topic
//
//	chatGptPromptList, err := QianfanServiceApp.parsePrompt(articleContext, ai.HotSpotWrite)
//	if err != nil {
//		return &ArticleContext{}, err
//	}
//
//
//	for _,
//	resp, err := chat.Do(
//		context.TODO(),
//		&qianfan.ChatCompletionRequest{
//			System: "微信公众号爆款文写作专家",
//			Messages: []qianfan.ChatCompletionMessage{
//				qianfan.ChatCompletionUserMessage(chatGptPrompt),
//			},
//		},
//	)
//
//	result := &ArticleContext{}
//
//	if err != nil {
//		global.GVA_LOG.Info("chat gpt响应失败", zap.Error(err))
//		return result, err
//	}
//
//	result.Content = resp.Result
//	articleContext.Content = resp.Result
//
//	chatGptPrompt, err = QianfanServiceApp.parsePrompt(articleContext, ai.TitleCreate)
//	if err != nil {
//		return nil, err
//	}
//
//	resp, err = chat.Do(
//		context.TODO(),
//		&qianfan.ChatCompletionRequest{
//			Messages: []qianfan.ChatCompletionMessage{
//				qianfan.ChatCompletionUserMessage(chatGptPrompt),
//			},
//		},
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	result.Title = resp.Result
//
//	chatGptPrompt, err = QianfanServiceApp.parsePrompt(articleContext, ai.AddImage)
//	if err != nil {
//		return nil, err
//	}
//
//	resp, err = chat.Do(
//		context.TODO(),
//		&qianfan.ChatCompletionRequest{
//			Messages: []qianfan.ChatCompletionMessage{
//				qianfan.ChatCompletionUserMessage(chatGptPrompt),
//			},
//		},
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	result.Content = resp.Result
//
//	return result, nil
//}
//
//func (*QianfanService) TopicWrite(topic string) (*ArticleContext, error) {
//	articleContext := &ArticleContext{}
//	articleContext.Topic = topic
//
//	chat := qianfan.NewChatCompletion(qianfan.WithModel("ERNIE-Bot-4"))
//
//	subject := SubjectServiceApp.FindAndUseSubjectByTopic(topic)
//	if subject == "" {
//		chatGptPrompt := "请以<" + topic + ">为主题随机提供一个有趣的吸引人的写作话题，直接返回话题即可，无需任何补充说明"
//		resp, err := chat.Do(
//			context.TODO(),
//			&qianfan.ChatCompletionRequest{
//				System: "微信公众号爆款文写作专家",
//				Messages: []qianfan.ChatCompletionMessage{
//					qianfan.ChatCompletionUserMessage(chatGptPrompt),
//				},
//			},
//		)
//		if err != nil {
//			return &ArticleContext{}, err
//		}
//		subject = resp.Result
//	}
//
//	articleContext.Topic = subject
//
//	chatGptPrompt, err := QianfanServiceApp.parsePrompt(articleContext, ai.TopicWrite)
//	if err != nil {
//		return &ArticleContext{}, err
//	}
//
//	resp, err := chat.Do(
//		context.TODO(),
//		&qianfan.ChatCompletionRequest{
//			System: "微信公众号爆款文写作专家",
//			Messages: []qianfan.ChatCompletionMessage{
//				qianfan.ChatCompletionUserMessage(chatGptPrompt),
//			},
//		},
//	)
//
//	if err != nil {
//		global.GVA_LOG.Info("chat gpt响应失败", zap.Error(err))
//		return nil, err
//	}
//
//	result := &ArticleContext{}
//	result.Content = resp.Result
//	articleContext.Content = resp.Result
//
//	chatGptPrompt, err = QianfanServiceApp.parsePrompt(articleContext, ai.TitleCreate)
//	if err != nil {
//		return articleContext, err
//	}
//
//	resp, err = chat.Do(
//		context.TODO(),
//		&qianfan.ChatCompletionRequest{
//			System: "微信公众号爆款文写作专家",
//			Messages: []qianfan.ChatCompletionMessage{
//				qianfan.ChatCompletionUserMessage(chatGptPrompt),
//			},
//		},
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	result.Title = resp.Result
//
//	chatGptPrompt, err = QianfanServiceApp.parsePrompt(articleContext, ai.AddImage)
//	if err != nil {
//		return nil, err
//	}
//
//	resp, err = chat.Do(
//		context.TODO(),
//		&qianfan.ChatCompletionRequest{
//			System: "微信公众号爆款文写作专家",
//			Messages: []qianfan.ChatCompletionMessage{
//				qianfan.ChatCompletionUserMessage(chatGptPrompt),
//			},
//		},
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	result.Content = resp.Result
//
//	return result, nil
//}
//
//func (*QianfanService) Recreation(article ai.Article) (*ArticleContext, error) {
//	articleContext := &ArticleContext{}
//	articleContext.Content = article.Content
//	articleContext.Title = article.Title
//	articleContext.Topic = article.Topic
//
//	// 可以通过 WithModel 指定模型
//	chat := qianfan.NewChatCompletion(qianfan.WithModel("ERNIE-Bot-4"))
//
//	chatGptPrompt, err := QianfanServiceApp.parsePrompt(articleContext, ai.ContentRecreation)
//	if err != nil {
//		return &ArticleContext{}, err
//	}
//
//	resp, err := chat.Do(
//		context.TODO(),
//		&qianfan.ChatCompletionRequest{
//			System: "微信公众号爆款文写作专家",
//			Messages: []qianfan.ChatCompletionMessage{
//				qianfan.ChatCompletionUserMessage(chatGptPrompt),
//			},
//		},
//	)
//
//	if err != nil {
//		global.GVA_LOG.Info("chat gpt响应失败", zap.Error(err))
//		return &ArticleContext{}, err
//	}
//
//	content := resp.Result
//	articleContext.Content = content
//
//	chatGptPrompt, err = QianfanServiceApp.parsePrompt(articleContext, ai.TitleCreate)
//	if err != nil {
//		return &ArticleContext{}, err
//	}
//	resp, err = chat.Do(
//		context.TODO(),
//		&qianfan.ChatCompletionRequest{
//			Messages: []qianfan.ChatCompletionMessage{
//				qianfan.ChatCompletionUserMessage(chatGptPrompt),
//			},
//		},
//	)
//
//	title := strings.TrimSpace(resp.Result)
//	title = strings.ReplaceAll(title, "{}", "")
//	articleContext.Title = title
//
//	chatGptPrompt, err = QianfanServiceApp.parsePrompt(articleContext, ai.AddImage)
//	if err != nil {
//		return nil, err
//	}
//
//	resp, err = chat.Do(
//		context.TODO(),
//		&qianfan.ChatCompletionRequest{
//			System: "微信公众号爆款文写作专家",
//			Messages: []qianfan.ChatCompletionMessage{
//				qianfan.ChatCompletionUserMessage(chatGptPrompt),
//			},
//		},
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	articleContext.Content = resp.Result
//
//	return articleContext, nil
//}
//

//
//type ChatGptResp struct {
//	Title   string
//	Content string
//	Topic   string
//	Tags    string
//}
