package ai

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
)

type QianwenService struct {
}

var QianwenServiceApp = new(QianwenService)

func (*QianwenService) GetKeyWord(title string) string {
	chatGptPrompt := "你现在是一名爆文写手，特别擅长从文章标题中找到关键词。我将给一个文章标题，需要你帮忙提取标题中的一个关键词用以做图片搜索。如果找不到关键词，可以返回该标题的主题，例如：历史，职场，明星等等" +
		"\n举例   " +
		"\n文章标题：中瑙友谊再升华，开启双边合作新篇章。  关键词：友谊再升华" +
		"\n文章标题：周处传奇：除三害、转人生，英雄之路的跌宕起伏  关键词：周处除三害" +
		"\n文章标题：" + title

	qianwenMessageHistory := []*QianwenMessage{QianwenSystemMessage}

	resp, qianwenMessageHistory, err := QianwenServiceApp.ChatWithQianwen(chatGptPrompt, qianwenMessageHistory)
	if err != nil || len(resp) > 10 {
		resp = "夜晚的城市"
	}

	return resp
}

func (*QianwenService) HotSpotWrite(context *ArticleContext) (*ArticleContext, error) {

	chatGptPromptList, err := parsePrompt(context, ai.ContentRecreation)
	if err != nil {
		return &ArticleContext{}, err
	}

	qianwenMessageHistory := []*QianwenMessage{QianwenSystemMessage}
	resp := ""

	for _, chatGptPrompt := range chatGptPromptList {
		resp, qianwenMessageHistory, err = QianwenServiceApp.ChatWithQianwen(chatGptPrompt, qianwenMessageHistory)
		if err != nil {
			return nil, err
		}
		context.Content = resp
	}

	chatGptPromptList, err = parsePrompt(context, ai.TitleCreate)
	if err != nil {
		return nil, err
	}

	for _, chatGptPrompt := range chatGptPromptList {
		resp, qianwenMessageHistory, err = QianwenServiceApp.ChatWithQianwen(chatGptPrompt, qianwenMessageHistory)
		if err != nil {
			return nil, err
		}
		context.Title = resp
	}

	chatGptPromptList, err = parsePrompt(context, ai.AddImage)
	if err != nil {
		return nil, err
	}

	for _, chatGptPrompt := range chatGptPromptList {
		resp, qianwenMessageHistory, err = QianwenServiceApp.ChatWithQianwen(chatGptPrompt, qianwenMessageHistory)
		if err != nil {
			return nil, err
		}

		context.Content = resp
	}

	context.Params = []string{"qianwen", "HotSpotWrite"}
	return context, nil
}

func (*QianwenService) TopicWrite(articleContext *ArticleContext) (*ArticleContext, error) {
	topic := articleContext.Topic

	subject := SubjectServiceApp.FindAndUseSubjectByTopic(topic)

	articleContext.Topic = subject

	chatGptPromptList, err := parsePrompt(articleContext, ai.TopicWrite)
	if err != nil {
		return &ArticleContext{}, err
	}

	qianwenMessageHistory := []*QianwenMessage{QianwenSystemMessage}
	resp := ""

	for _, chatGptPrompt := range chatGptPromptList {
		resp, qianwenMessageHistory, err = QianwenServiceApp.ChatWithQianwen(chatGptPrompt, qianwenMessageHistory)
		if err != nil {
			return nil, err
		}
		articleContext.Content = resp
	}

	chatGptPromptList, err = parsePrompt(articleContext, ai.TitleCreate)
	if err != nil {
		return nil, err
	}

	for _, chatGptPrompt := range chatGptPromptList {
		resp, qianwenMessageHistory, err = QianwenServiceApp.ChatWithQianwen(chatGptPrompt, qianwenMessageHistory)
		if err != nil {
			return nil, err
		}
		articleContext.Title = resp
	}

	chatGptPromptList, err = parsePrompt(articleContext, ai.AddImage)
	if err != nil {
		return nil, err
	}

	for _, chatGptPrompt := range chatGptPromptList {
		resp, qianwenMessageHistory, err = QianwenServiceApp.ChatWithQianwen(chatGptPrompt, qianwenMessageHistory)
		if err != nil {
			return nil, err
		}

		articleContext.Content = resp
	}

	articleContext.Params = []string{"qianwen", "HotSpotWrite"}

	return articleContext, nil
}

func (*QianwenService) Recreation(articleContext *ArticleContext) (*ArticleContext, error) {

	chatGptPromptList, err := parsePrompt(articleContext, ai.ContentRecreation)
	if err != nil {
		return &ArticleContext{}, err
	}

	qianwenMessageHistory := []*QianwenMessage{QianwenSystemMessage}
	resp := ""

	for _, chatGptPrompt := range chatGptPromptList {
		resp, qianwenMessageHistory, err = QianwenServiceApp.ChatWithQianwen(chatGptPrompt, qianwenMessageHistory)
		if err != nil {
			return nil, err
		}
		articleContext.Content = resp
	}

	chatGptPromptList, err = parsePrompt(articleContext, ai.TitleCreate)
	if err != nil {
		return nil, err
	}

	for _, chatGptPrompt := range chatGptPromptList {
		resp, qianwenMessageHistory, err = QianwenServiceApp.ChatWithQianwen(chatGptPrompt, qianwenMessageHistory)
		if err != nil {
			return nil, err
		}
		articleContext.Title = resp
	}

	chatGptPromptList, err = parsePrompt(articleContext, ai.AddImage)
	if err != nil {
		return nil, err
	}

	for _, chatGptPrompt := range chatGptPromptList {
		resp, qianwenMessageHistory, err = QianwenServiceApp.ChatWithQianwen(chatGptPrompt, qianwenMessageHistory)
		if err != nil {
			return nil, err
		}

		articleContext.Content = resp
	}

	articleContext.Params = []string{"qianwen", "Recreation"}
	return articleContext, nil
}

var QianwenSystemMessage = &QianwenMessage{
	Role:    "system",
	Content: "你是 Qianwen，是一个人工智能助手，你更擅长中文和英文的对话。你会为用户提供安全，有帮助，准确的回答。同时，你会拒绝一切涉及恐怖主义，种族歧视，黄色暴力等问题的回答。Moonshot AI 为专有名词，不可翻译成其他语言。",
}

func (*QianwenService) ChatWithQianwen(message string, history []*QianwenMessage) (string, []*QianwenMessage, error) {
	kimiCfg := global.GVA_CONFIG.Qianwen

	history = append(history, &QianwenMessage{
		Role:    "user",
		Content: message,
	})

	kimiReq := &QianwenReq{
		Model:    "moonshot-v1-8k",
		Messages: history,
	}

	apiUrl := kimiCfg.ApiUrl

	statusCode, respBody, err := utils.PostWithHeaders(apiUrl, utils.Parse2Json(kimiReq), map[string]string{
		"Authorization": "Bearer " + kimiCfg.RefreshToken,
	})

	if err != nil {
		return "", history, err
	}

	if statusCode != 200 {
		return "", history, errors.New(string(respBody))
	}

	kimiResp := &QianwenResp{}
	err = utils.JsonStrToStruct(string(respBody), kimiResp)
	if err != nil {
		return "", history, err
	}

	if len(kimiResp.Choices) > 0 {
		history = append(history, &QianwenMessage{
			Role:    "assistant",
			Content: kimiResp.Choices[0].Message.Content,
		})
		return kimiResp.Choices[0].Message.Content, history, err
	}

	return "", history, errors.New("kimi回复为空")
}

type QianwenReq struct {
	Model       string            `json:"model"`
	Messages    []*QianwenMessage `json:"messages"`
	Temperature float64           `json:"temperature"`
}

type QianwenMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type QianwenResp struct {
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
