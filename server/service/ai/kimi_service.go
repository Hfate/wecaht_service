package ai

import (
	"bytes"
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"go.uber.org/zap"
	"strings"
	"text/template"
)

type KimiService struct {
}

var KimiServiceApp = new(KimiService)

func (*KimiService) HotSpotWrite(topic string) (*ArticleContext, error) {

	chatGptPrompt := "请以<" + topic + ">为话题，结合你的联网搜索能力，写一篇1200字的文章，文章内容各处无需补充说明,要求返回文章格式为markdown格式，且限定markdown仅支持字体加粗，下划线，斜体，有序列表等格式"
	resp, err := KimiServiceApp.ChatWithKimi(chatGptPrompt)

	result := &ArticleContext{}

	if err != nil {
		global.GVA_LOG.Info("chat gpt响应失败", zap.Error(err))
		return result, err
	}

	result.Content = resp

	chatGptPrompt = "请给下文生成一个吸引人阅读的标题，直接回答标题即可无需补充说明。文章内容:" + result.Content

	resp, err = KimiServiceApp.ChatWithKimi(chatGptPrompt)

	result.Title = resp

	chatGptPrompt = "你是一位微信公众号爆文写手，擅长创作爆款文章并为其配上合适的图片。现在，你需要基于提供的文章内容，在合适的位置添加配图占位符，以提升读者的阅读体验。" +
		" 占位符的格式示例：[img]高中生放学[/img]。" +
		" 请确保图片与文章主题和内容紧密相连，让读者在阅读过程中能够更好地理解和感受文章所传达的信息和情感。" +
		" 原文如下:" + result.Content

	resp, err = KimiServiceApp.ChatWithKimi(chatGptPrompt)
	result.Content = resp

	return result, nil
}

func (*KimiService) TopicWrite(topic string) (*ArticleContext, error) {

	chatGptPrompt := "请以<" + topic + ">为主题，结合你的联网搜索能力，随机提供一个有趣的不重复的紧贴时事的写作话题，直接返回话题即可"

	resp, err := KimiServiceApp.ChatWithKimi(chatGptPrompt)

	topic = resp

	chatGptPrompt = "请以<" + topic + ">为主题写一篇1200字的微信公众号文章，文章内容各处无需补充说明，要求返回文章格式为markdown格式，且限定markdown仅支持字体加粗，下划线，斜体，有序列表等格式"
	resp, err = KimiServiceApp.ChatWithKimi(chatGptPrompt)

	result := &ArticleContext{}

	if err != nil {
		global.GVA_LOG.Info("chat gpt响应失败", zap.Error(err))
		return result, err
	}

	result.Content = resp

	chatGptPrompt = "请给下文生成一个吸引人阅读的标题，直接回答标题即可无需补充说明.文章内容" + result.Content

	resp, err = KimiServiceApp.ChatWithKimi(chatGptPrompt)

	result.Title = resp

	chatGptPrompt = "你是一位微信公众号爆文写手，擅长创作爆款文章并为其配上合适的图片。现在，你需要基于提供的文章内容，在合适的位置添加配图占位符，以提升读者的阅读体验。" +
		" 占位符的格式示例：[img]高中生放学[/img]。" +
		" 1 请确保图片与文章主题和内容紧密相连，让读者在阅读过程中能够更好地理解和感受文章所传达的信息和情感。" +
		" 2 直接返回增加配图占位符后的文章即可，无需任何补充说明" +
		" 原文如下:" + result.Content

	resp, err = KimiServiceApp.ChatWithKimi(chatGptPrompt)
	result.Content = resp

	return result, nil
}

func (*KimiService) Recreation(article ai.Article) (*ArticleContext, error) {

	chatGptPrompt := KimiServiceApp.parseContentPrompt(article)

	resp, err := KimiServiceApp.ChatWithKimi(chatGptPrompt)

	if err != nil {
		global.GVA_LOG.Info("chat gpt响应失败", zap.Error(err))
		return &ArticleContext{}, err
	}

	content := resp

	chatGptResp := &ArticleContext{}
	chatGptResp.Content = content

	chatGptPrompt = KimiServiceApp.parseTitlePrompt(article)
	resp, err = KimiServiceApp.ChatWithKimi(chatGptPrompt)

	title := strings.TrimSpace(resp)
	title = strings.ReplaceAll(title, "{}", "")
	chatGptResp.Title = title

	return chatGptResp, nil
}

func (*KimiService) parseContentPrompt(article ai.Article) string {
	topic := article.Topic
	prompt, err := PromptServiceApp.FindPromptByTopicAndType(topic, ai.ContentRecreation)
	if err != nil {
		global.GVA_LOG.Info("无法找到topic相关的prompt", zap.Error(err), zap.String("topic", topic))
		return ""
	}

	temp := template.New("ChatGptPrompt")
	tmpl, err := temp.Parse(prompt)
	if err != nil {
		global.GVA_LOG.Info("模板解析失败")
		return ""
	}

	// 创建一个缓冲区来保存模板生成的结果
	var buf bytes.Buffer
	// 使用模板和数据生成输出
	err = tmpl.Execute(&buf, article)

	chatGptPrompt := buf.String()
	return chatGptPrompt
}

func (*KimiService) parseTitlePrompt(article ai.Article) string {
	topic := article.Topic
	prompt, err := PromptServiceApp.FindPromptByTopicAndType(topic, ai.TitleRecreation)
	if err != nil {
		global.GVA_LOG.Info("无法找到topic相关的prompt", zap.Error(err), zap.String("topic", topic))
		return ""
	}

	temp := template.New("ChatGptPrompt")
	tmpl, err := temp.Parse(prompt)
	if err != nil {
		global.GVA_LOG.Info("模板解析失败")
		return ""
	}

	// 创建一个缓冲区来保存模板生成的结果
	var buf bytes.Buffer
	// 使用模板和数据生成输出
	err = tmpl.Execute(&buf, article)

	chatGptPrompt := buf.String()
	return chatGptPrompt
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
