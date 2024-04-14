package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/baidubce/bce-qianfan-sdk/go/qianfan"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"text/template"
)

type KimiService struct {
}

var KimiServiceApp = new(KimiService)

func (*KimiService) GetKeyWord(title string) string {
	chatGptPrompt := "你现在是一名爆文写手，特别擅长从文章标题中找到关键词。我将给一个文章标题，需要你帮忙提取标题中的一个关键词用以做图片搜索。如果找不到关键词，可以返回该标题的主题，例如：历史，职场，明星等等" +
		"\n举例   " +
		"\n文章标题：中瑙友谊再升华，开启双边合作新篇章。  关键词：友谊再升华" +
		"\n文章标题：周处传奇：除三害、转人生，英雄之路的跌宕起伏  关键词：周处除三害" +
		"\n文章标题：" + title

	kimiMessageHistory := []*KimiMessage{SystemMessage}

	resp, kimiMessageHistory, err := KimiServiceApp.ChatWithKimi(chatGptPrompt, kimiMessageHistory)
	if err != nil || len(resp) > 10 {
		resp = "夜晚的城市"
	}

	return resp
}

func (*KimiService) HotSpotWrite(context *ArticleContext) (*ArticleContext, error) {

	chatGptPromptList, err := parsePrompt(context, ai.ContentRecreation)
	if err != nil {
		return &ArticleContext{}, err
	}

	kimiMessageHistory := []*KimiMessage{SystemMessage}
	result := &ArticleContext{}
	resp := ""

	for _, chatGptPrompt := range chatGptPromptList {
		resp, kimiMessageHistory, err = KimiServiceApp.ChatWithKimi(chatGptPrompt, kimiMessageHistory)
		if err != nil {
			return nil, err
		}
		result.Content = resp
		context.Content = resp
	}

	chatGptPromptList, err = parsePrompt(context, ai.TitleCreate)
	if err != nil {
		return nil, err
	}

	for _, chatGptPrompt := range chatGptPromptList {
		resp, kimiMessageHistory, err = KimiServiceApp.ChatWithKimi(chatGptPrompt, kimiMessageHistory)
		if err != nil {
			return nil, err
		}
		result.Title = resp
	}

	chatGptPromptList, err = parsePrompt(context, ai.AddImage)
	if err != nil {
		return nil, err
	}

	for _, chatGptPrompt := range chatGptPromptList {
		resp, kimiMessageHistory, err = KimiServiceApp.ChatWithKimi(chatGptPrompt, kimiMessageHistory)
		if err != nil {
			return nil, err
		}

		result.Content = resp
	}

	context.Params = []string{"kimi", "HotWrite"}
	return result, nil
}

func (*KimiService) TopicWrite(articleContext *ArticleContext) (*ArticleContext, error) {
	topic := articleContext.Topic

	subject := SubjectServiceApp.FindAndUseSubjectByTopic(topic)
	if subject == "" {
		chat := qianfan.NewChatCompletion(qianfan.WithModel("ERNIE-Bot-4"))
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

	chatGptPromptList, err := parsePrompt(articleContext, ai.TopicWrite)
	if err != nil {
		return &ArticleContext{}, err
	}

	kimiMessageHistory := []*KimiMessage{SystemMessage}
	result := &ArticleContext{}
	resp := ""

	for _, chatGptPrompt := range chatGptPromptList {
		resp, kimiMessageHistory, err = KimiServiceApp.ChatWithKimi(chatGptPrompt, kimiMessageHistory)
		if err != nil {
			return nil, err
		}
		result.Content = resp
		articleContext.Content = resp
	}

	chatGptPromptList, err = parsePrompt(articleContext, ai.TitleCreate)
	if err != nil {
		return nil, err
	}

	for _, chatGptPrompt := range chatGptPromptList {
		resp, kimiMessageHistory, err = KimiServiceApp.ChatWithKimi(chatGptPrompt, kimiMessageHistory)
		if err != nil {
			return nil, err
		}
		result.Title = resp
	}

	chatGptPromptList, err = parsePrompt(articleContext, ai.AddImage)
	if err != nil {
		return nil, err
	}

	for _, chatGptPrompt := range chatGptPromptList {
		resp, kimiMessageHistory, err = KimiServiceApp.ChatWithKimi(chatGptPrompt, kimiMessageHistory)
		if err != nil {
			return nil, err
		}

		result.Content = resp
	}

	articleContext.Params = []string{"kimi", "TopicWrite"}
	return result, nil
}

func (*KimiService) Recreation(context *ArticleContext) (*ArticleContext, error) {

	chatGptPromptList, err := parsePrompt(context, ai.ContentRecreation)
	if err != nil {
		return &ArticleContext{}, err
	}

	kimiMessageHistory := []*KimiMessage{SystemMessage}

	resp := ""

	for _, chatGptPrompt := range chatGptPromptList {
		resp, kimiMessageHistory, err = KimiServiceApp.ChatWithKimi(chatGptPrompt, kimiMessageHistory)
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
		resp, kimiMessageHistory, err = KimiServiceApp.ChatWithKimi(chatGptPrompt, kimiMessageHistory)
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
		resp, kimiMessageHistory, err = KimiServiceApp.ChatWithKimi(chatGptPrompt, kimiMessageHistory)
		if err != nil {
			return nil, err
		}

		context.Content = resp
	}

	context.Params = []string{"kimi", "Recreation"}
	return context, nil
}

var SystemMessage = &KimiMessage{
	Role:    "system",
	Content: "你是 Kimi，由 Moonshot AI 提供的人工智能助手，你更擅长中文和英文的对话。你会为用户提供安全，有帮助，准确的回答。同时，你会拒绝一切涉及恐怖主义，种族歧视，黄色暴力等问题的回答。Moonshot AI 为专有名词，不可翻译成其他语言。",
}

func (*KimiService) ChatWithKimi(message string, history []*KimiMessage) (string, []*KimiMessage, error) {
	kimiCfg := global.GVA_CONFIG.Kimi

	history = append(history, &KimiMessage{
		Role:    "user",
		Content: message,
	})

	kimiReq := &KimiReq{
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

	kimiResp := &KimiResp{}
	err = utils.JsonStrToStruct(string(respBody), kimiResp)
	if err != nil {
		return "", history, err
	}

	if len(kimiResp.Choices) > 0 {
		history = append(history, &KimiMessage{
			Role:    "assistant",
			Content: kimiResp.Choices[0].Message.Content,
		})
		return kimiResp.Choices[0].Message.Content, history, err
	}

	return "", history, errors.New("kimi回复为空")
}

type KimiReq struct {
	Model       string         `json:"model"`
	Messages    []*KimiMessage `json:"messages"`
	Temperature float64        `json:"temperature"`
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

func parsePrompt(context *ArticleContext, promptType int) ([]string, error) {
	topic := context.Topic
	prompt, err := PromptServiceApp.FindPromptByTopicAndType(topic, promptType)
	if err != nil {
		return []string{}, err
	}

	promptList := make([]string, 0)
	// 使用json.Unmarshal将JSON字符串解析到字符串切片
	err = json.Unmarshal([]byte(utils.EscapeSpecialCharacters(prompt)), &promptList)

	if len(promptList) == 0 {
		promptList = []string{prompt}
	}

	result := make([]string, 0)

	for _, item := range promptList {
		temp := template.New("ChatGptPrompt")
		tmpl, err2 := temp.Parse(item)
		if err2 != nil {
			global.GVA_LOG.Info("模板解析失败")
			return []string{}, err2
		}

		// 创建一个缓冲区来保存模板生成的结果
		var buf bytes.Buffer
		// 使用模板和数据生成输出
		err = tmpl.Execute(&buf, context)

		chatGptPrompt := buf.String()

		result = append(result, chatGptPrompt)
	}

	return result, err
}
