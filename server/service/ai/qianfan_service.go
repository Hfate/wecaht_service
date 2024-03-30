package ai

import (
	"bytes"
	"context"
	"github.com/baidubce/bce-qianfan-sdk/go/qianfan"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"go.uber.org/zap"
	"strings"
	"text/template"
)

type QianfanService struct {
}

var QianfanServiceApp = new(QianfanService)

func (*QianfanService) GetKeyWord(title string) string {
	chat := qianfan.NewChatCompletion(qianfan.WithModel("ERNIE-Bot-4"))

	chatGptPrompt := "我将给一个文章标题，需要你帮忙提取标题中的一个关键词用以做图片搜索。" +
		"\n举例   " +
		"\n文章标题：中瑙友谊再升华，开启双边合作新篇章。  关键词：友谊再升华" +
		"\n文章标题：周处传奇：除三害、转人生，英雄之路的跌宕起伏  关键词：周处除三害" +
		"\n文章标题：" + title
	resp, _ := chat.Do(
		context.TODO(),
		&qianfan.ChatCompletionRequest{
			System: "微信公众号爆款文写作专家",
			Messages: []qianfan.ChatCompletionMessage{
				qianfan.ChatCompletionUserMessage(chatGptPrompt),
			},
		},
	)
	return resp.Result
}

func (*QianfanService) HotSpotWrite(topic string) (*ArticleContext, error) {
	chat := qianfan.NewChatCompletion(qianfan.WithModel("ERNIE-Bot-4"))

	chatGptPrompt := "请以<" + topic + ">为主题写一篇1200字的文章，文章内容各处无需补充说明,要求返回文章格式为markdown格式"
	resp, err := chat.Do(
		context.TODO(),
		&qianfan.ChatCompletionRequest{
			System: "微信公众号爆款文写作专家",
			Messages: []qianfan.ChatCompletionMessage{
				qianfan.ChatCompletionUserMessage(chatGptPrompt),
			},
		},
	)

	result := &ArticleContext{}

	if err != nil {
		global.GVA_LOG.Info("chat gpt响应失败", zap.Error(err))
		return result, err
	}

	result.Content = resp.Result

	chatGptPrompt = "请给下文生成一个吸引人阅读的标题，直接回答标题即可无需补充说明。文章内容:" + result.Content

	resp, err = chat.Do(
		context.TODO(),
		&qianfan.ChatCompletionRequest{
			Messages: []qianfan.ChatCompletionMessage{
				qianfan.ChatCompletionUserMessage(chatGptPrompt),
			},
		},
	)

	result.Title = resp.Result

	chatGptPrompt = "你是一位微信公众号爆文写手，擅长创作爆款文章并为其配上合适的图片。现在，你需要基于提供的文章内容，在合适的位置添加配图占位符，以提升读者的阅读体验。" +
		" 占位符的格式示例：[img]高中生放学[/img]。" +
		" 请确保图片与文章主题和内容紧密相连，让读者在阅读过程中能够更好地理解和感受文章所传达的信息和情感。" +
		" 原文如下:" + result.Content

	resp, err = chat.Do(
		context.TODO(),
		&qianfan.ChatCompletionRequest{
			Messages: []qianfan.ChatCompletionMessage{
				qianfan.ChatCompletionUserMessage(chatGptPrompt),
			},
		},
	)
	result.Content = resp.Result

	return result, nil
}

func (*QianfanService) TopicSpotWrite(topic string) (*ArticleContext, error) {

	chat := qianfan.NewChatCompletion(qianfan.WithModel("ERNIE-Bot-4"))

	chatGptPrompt := "请以<" + topic + ">为主题随机提供一个有趣的不重复紧贴时事的写作话题，直接返回即可"

	resp, err := chat.Do(
		context.TODO(),
		&qianfan.ChatCompletionRequest{
			System: "微信公众号爆款文写作专家",
			Messages: []qianfan.ChatCompletionMessage{
				qianfan.ChatCompletionUserMessage(chatGptPrompt),
			},
		},
	)

	topic = resp.Result

	chatGptPrompt = "请以<" + topic + ">为主题写一篇1200字的微信公众号文章，文章内容各处无需补充说明，要求返回文章格式为markdown格式"
	resp, err = chat.Do(
		context.TODO(),
		&qianfan.ChatCompletionRequest{
			System: "微信公众号爆款文写作专家",
			Messages: []qianfan.ChatCompletionMessage{
				qianfan.ChatCompletionUserMessage(chatGptPrompt),
			},
		},
	)

	result := &ArticleContext{}

	if err != nil {
		global.GVA_LOG.Info("chat gpt响应失败", zap.Error(err))
		return result, err
	}

	result.Content = resp.Result

	chatGptPrompt = "请给下文生成一个吸引人阅读的标题，直接回答标题即可无需补充说明.文章内容" + result.Content

	resp, err = chat.Do(
		context.TODO(),
		&qianfan.ChatCompletionRequest{
			System: "微信公众号爆款文写作专家",
			Messages: []qianfan.ChatCompletionMessage{
				qianfan.ChatCompletionUserMessage(chatGptPrompt),
			},
		},
	)

	result.Title = resp.Result

	chatGptPrompt = "你是一位微信公众号爆文写手，擅长创作爆款文章并为其配上合适的图片。现在，你需要基于提供的文章内容，在合适的位置添加配图占位符，以提升读者的阅读体验。" +
		" 占位符的格式示例：[img]高中生放学[/img]。" +
		" 请确保图片与文章主题和内容紧密相连，让读者在阅读过程中能够更好地理解和感受文章所传达的信息和情感。" +
		" 原文如下:" + result.Content

	resp, err = chat.Do(
		context.TODO(),
		&qianfan.ChatCompletionRequest{
			System: "微信公众号爆款文写作专家",
			Messages: []qianfan.ChatCompletionMessage{
				qianfan.ChatCompletionUserMessage(chatGptPrompt),
			},
		},
	)
	result.Content = resp.Result

	return result, nil
}

func (*QianfanService) Recreation(article ai.Article) (*ArticleContext, error) {
	// 可以通过 WithModel 指定模型
	chat := qianfan.NewChatCompletion(qianfan.WithModel("ERNIE-Bot-4"))

	chatGptPrompt := QianfanServiceApp.parseContentPrompt(article)

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
		global.GVA_LOG.Info("chat gpt响应失败", zap.Error(err))
		return &ArticleContext{}, err
	}

	content := resp.Result

	chatGptResp := &ArticleContext{}
	chatGptResp.Content = content

	chatGptPrompt = QianfanServiceApp.parseTitlePrompt(article)
	resp, err = chat.Do(
		context.TODO(),
		&qianfan.ChatCompletionRequest{
			Messages: []qianfan.ChatCompletionMessage{
				qianfan.ChatCompletionUserMessage(chatGptPrompt),
			},
		},
	)

	title := strings.TrimSpace(resp.Result)
	title = strings.ReplaceAll(title, "{}", "")
	chatGptResp.Title = title

	return chatGptResp, nil
}

func (*QianfanService) parseContentPrompt(article ai.Article) string {
	topic := article.Topic
	prompt, err := PromptServiceApp.FindPromptByTopicAndType(topic, 1)
	// 没找到 则使用默认的
	if err != nil {
		prompt, err = PromptServiceApp.FindPromptByTopicAndType("default", 1)
		if err != nil {
			global.GVA_LOG.Info("无法找到topic相关的prompt", zap.Error(err), zap.String("topic", topic))
			return ""
		}
	}

	temp := template.New("ChatGptPrompt")
	tmpl, err := temp.Parse(prompt.Prompt)
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

func (*QianfanService) parseTitlePrompt(article ai.Article) string {
	topic := article.Topic
	prompt, err := PromptServiceApp.FindPromptByTopicAndType(topic, 2)
	if err != nil {
		prompt, err = PromptServiceApp.FindPromptByTopicAndType("default", 2)
		if err != nil {
			global.GVA_LOG.Info("无法找到topic相关的prompt", zap.Error(err), zap.String("topic", topic))
			return ""
		}
	}

	temp := template.New("ChatGptPrompt")
	tmpl, err := temp.Parse(prompt.Prompt)
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

type ChatGptResp struct {
	Title   string
	Content string
	Topic   string
	Tags    string
}
