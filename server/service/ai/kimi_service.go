package ai

import (
	"context"
	"errors"
	"github.com/baidubce/bce-qianfan-sdk/go/qianfan"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"go.uber.org/zap"
	"strings"
)

type KimiService struct {
}

var KimiServiceApp = new(KimiService)

func (*KimiService) HotSpotWrite(topic string) (*ArticleContext, error) {

	articleContext := &ArticleContext{}
	articleContext.Topic = topic

	chatGptPrompt := "请以<" + topic + ">为话题，结合你的联网搜索能力，写一篇1200字的文章，文章内容各处无需补充说明,要求返回文章格式为markdown格式，且限定markdown仅支持字体加粗，下划线，斜体，有序列表等格式"
	resp, err := KimiServiceApp.ChatWithKimi(chatGptPrompt)

	result := &ArticleContext{}

	if err != nil {
		global.GVA_LOG.Info("chat gpt响应失败", zap.Error(err))
		return result, err
	}

	result.Content = resp
	articleContext.Content = resp

	chatGptPrompt, err = QianfanServiceApp.parsePrompt(articleContext, ai.TitleCreate)
	if err != nil {
		return nil, err
	}

	resp, err = KimiServiceApp.ChatWithKimi(chatGptPrompt)
	if err != nil {
		return nil, err
	}

	result.Title = resp

	chatGptPrompt, err = QianfanServiceApp.parsePrompt(articleContext, ai.AddImage)
	if err != nil {
		return nil, err
	}

	resp, err = KimiServiceApp.ChatWithKimi(chatGptPrompt)
	if err != nil {
		return nil, err
	}

	result.Content = resp

	return result, nil
}

func (*KimiService) TopicWrite(topic string) (*ArticleContext, error) {
	articleContext := &ArticleContext{}
	articleContext.Topic = topic

	chat := qianfan.NewChatCompletion(qianfan.WithModel("ERNIE-Bot-4"))

	subject := SubjectServiceApp.FindAndUseSubjectByTopic(topic)
	if subject == "" {
		chatGptPrompt := "请以<" + topic + ">为主题随机提供一个有趣的吸引人的写作话题，直接返回话题即可，无需任何补充说明"
		resp, err := chat.Do(
			context.TODO(),
			&qianfan.ChatCompletionRequest{
				System: "微信公众号爆款文写作专家",
				Messages: []qianfan.ChatCompletionMessage{
					qianfan.ChatCompletionUserMessage(chatGptPrompt),
				},
			},
		)
		if err != nil {
			return &ArticleContext{}, err
		}
		subject = resp.Result
	}

	articleContext.Topic = subject

	chatGptPrompt, err := QianfanServiceApp.parsePrompt(articleContext, ai.TopicWrite)
	if err != nil {
		return &ArticleContext{}, err
	}

	resp, err := KimiServiceApp.ChatWithKimi(chatGptPrompt)

	if err != nil {
		global.GVA_LOG.Info("chat gpt响应失败", zap.Error(err))
		return nil, err
	}

	result := &ArticleContext{}
	result.Content = resp
	articleContext.Content = resp

	chatGptPrompt, err = QianfanServiceApp.parsePrompt(articleContext, ai.TitleCreate)
	if err != nil {
		return articleContext, err
	}

	resp, err = KimiServiceApp.ChatWithKimi(chatGptPrompt)
	if err != nil {
		return nil, err
	}

	result.Title = resp

	chatGptPrompt, err = QianfanServiceApp.parsePrompt(articleContext, ai.AddImage)
	if err != nil {
		return nil, err
	}

	resp, err = KimiServiceApp.ChatWithKimi(chatGptPrompt)
	if err != nil {
		return nil, err
	}

	result.Content = resp

	return result, nil
}

func (*KimiService) Recreation(article ai.Article) (*ArticleContext, error) {
	articleContext := &ArticleContext{}
	articleContext.Content = article.Content
	articleContext.Title = article.Title
	articleContext.Topic = article.Topic

	chatGptPrompt, err := QianfanServiceApp.parsePrompt(articleContext, ai.ContentRecreation)
	if err != nil {
		return &ArticleContext{}, err
	}

	resp, err := KimiServiceApp.ChatWithKimi(chatGptPrompt)

	if err != nil {
		global.GVA_LOG.Info("chat gpt响应失败", zap.Error(err))
		return &ArticleContext{}, err
	}

	content := resp
	articleContext.Content = content

	chatGptPrompt, err = QianfanServiceApp.parsePrompt(articleContext, ai.TitleCreate)
	if err != nil {
		return &ArticleContext{}, err
	}
	resp, err = KimiServiceApp.ChatWithKimi(chatGptPrompt)

	title := strings.TrimSpace(resp)
	title = strings.ReplaceAll(title, "{}", "")
	articleContext.Title = title

	chatGptPrompt, err = QianfanServiceApp.parsePrompt(articleContext, ai.AddImage)
	if err != nil {
		return nil, err
	}

	resp, err = KimiServiceApp.ChatWithKimi(chatGptPrompt)
	if err != nil {
		return nil, err
	}

	articleContext.Content = resp

	return articleContext, nil
}

func (*KimiService) ChatWithKimi(message string) (string, error) {
	kimiCfg := global.GVA_CONFIG.Kimi

	kimiReq := &KimiReq{
		Model: "moonshot-v1-8k",
		Messages: []KimiMessage{{
			Role:    "system",
			Content: "你是 Kimi，由 Moonshot AI 提供的人工智能助手，你更擅长中文和英文的对话。你会为用户提供安全，有帮助，准确的回答。同时，你会拒绝一切涉及恐怖主义，种族歧视，黄色暴力等问题的回答。Moonshot AI 为专有名词，不可翻译成其他语言。",
		}, {
			Role:    "user",
			Content: message,
		},
		},
		Temperature: 0.3,
	}

	statusCode, respBody, err := utils.PostWithHeaders("https://api.moonshot.cn/v1/chat/completions", utils.Parse2Json(kimiReq), map[string]string{
		"Authorization": "Bearer " + kimiCfg.AccessKey,
	})

	if err != nil {
		return "", err
	}

	if statusCode != 200 {
		return "", errors.New(string(respBody))
	}

	kimiResp := &KimiResp{}
	err = utils.JsonStrToStruct(string(respBody), kimiResp)
	if err != nil {
		return "", err
	}

	return kimiResp.Choices[0].Message.Content, err
}

type KimiReq struct {
	Model       string        `json:"model"`
	Messages    []KimiMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
}

type KimiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type KimiResp struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}
