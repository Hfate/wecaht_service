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

	chatGptPrompt := "我将给一个文章标题，需要你帮忙提取标题中的一个关键词用以做图片搜索。如果找不到关键词，可以返回该标题的主题，例如：历史，职场，明星等等" +
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
	articleContext := &ArticleContext{}
	articleContext.Topic = topic

	chatGptPrompt, err := QianfanServiceApp.parsePrompt(articleContext, ai.HotSpotWrite)
	if err != nil {
		return &ArticleContext{}, err
	}

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
	articleContext.Content = resp.Result

	chatGptPrompt, err = QianfanServiceApp.parsePrompt(articleContext, ai.TitleCreate)
	if err != nil {
		return nil, err
	}

	resp, err = chat.Do(
		context.TODO(),
		&qianfan.ChatCompletionRequest{
			Messages: []qianfan.ChatCompletionMessage{
				qianfan.ChatCompletionUserMessage(chatGptPrompt),
			},
		},
	)
	if err != nil {
		return nil, err
	}

	result.Title = resp.Result

	chatGptPrompt, err = QianfanServiceApp.parsePrompt(articleContext, ai.AddImage)
	if err != nil {
		return nil, err
	}

	resp, err = chat.Do(
		context.TODO(),
		&qianfan.ChatCompletionRequest{
			Messages: []qianfan.ChatCompletionMessage{
				qianfan.ChatCompletionUserMessage(chatGptPrompt),
			},
		},
	)
	if err != nil {
		return nil, err
	}

	result.Content = resp.Result

	return result, nil
}

func (*QianfanService) TopicWrite(topic string) (*ArticleContext, error) {
	articleContext := &ArticleContext{}
	articleContext.Topic = topic

	chat := qianfan.NewChatCompletion(qianfan.WithModel("ERNIE-Bot-4"))

	chatGptPrompt := "请以<" + topic + ">为主题随机提供一个有趣的不重复紧贴时事的写作话题，直接返回话题即可，无需任何补充说明"

	resp, err := chat.Do(
		context.TODO(),
		&qianfan.ChatCompletionRequest{
			System: "微信公众号爆款文写作专家",
			Messages: []qianfan.ChatCompletionMessage{
				qianfan.ChatCompletionUserMessage(chatGptPrompt),
			},
		},
	)

	articleContext.Topic = resp.Result

	chatGptPrompt, err = QianfanServiceApp.parsePrompt(articleContext, ai.TopicWrite)
	if err != nil {
		return &ArticleContext{}, err
	}

	resp, err = chat.Do(
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
		return nil, err
	}

	result := &ArticleContext{}
	result.Content = resp.Result
	articleContext.Content = resp.Result

	chatGptPrompt, err = QianfanServiceApp.parsePrompt(articleContext, ai.TitleCreate)
	if err != nil {
		return articleContext, err
	}

	resp, err = chat.Do(
		context.TODO(),
		&qianfan.ChatCompletionRequest{
			System: "微信公众号爆款文写作专家",
			Messages: []qianfan.ChatCompletionMessage{
				qianfan.ChatCompletionUserMessage(chatGptPrompt),
			},
		},
	)
	if err != nil {
		return nil, err
	}

	result.Title = resp.Result

	chatGptPrompt, err = QianfanServiceApp.parsePrompt(articleContext, ai.AddImage)
	if err != nil {
		return nil, err
	}

	resp, err = chat.Do(
		context.TODO(),
		&qianfan.ChatCompletionRequest{
			System: "微信公众号爆款文写作专家",
			Messages: []qianfan.ChatCompletionMessage{
				qianfan.ChatCompletionUserMessage(chatGptPrompt),
			},
		},
	)
	if err != nil {
		return nil, err
	}

	result.Content = resp.Result

	return result, nil
}

func (*QianfanService) Recreation(article ai.Article) (*ArticleContext, error) {
	articleContext := &ArticleContext{}
	articleContext.Content = article.Content
	articleContext.Title = article.Title
	articleContext.Topic = article.Topic

	// 可以通过 WithModel 指定模型
	chat := qianfan.NewChatCompletion(qianfan.WithModel("ERNIE-Bot-4"))

	chatGptPrompt, err := QianfanServiceApp.parsePrompt(articleContext, ai.ContentRecreation)
	if err != nil {
		return &ArticleContext{}, err
	}

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
	articleContext.Content = content

	chatGptPrompt, err = QianfanServiceApp.parsePrompt(articleContext, ai.TitleCreate)
	if err != nil {
		return &ArticleContext{}, err
	}
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
	articleContext.Title = title

	return articleContext, nil
}

func (*QianfanService) parsePrompt(context *ArticleContext, promptType int) (string, error) {
	topic := context.Topic
	prompt, err := PromptServiceApp.FindPromptByTopicAndType(topic, promptType)
	if err != nil {
		return "", err
	}

	temp := template.New("ChatGptPrompt")
	tmpl, err := temp.Parse(prompt)
	if err != nil {
		global.GVA_LOG.Info("模板解析失败")
		return "", err
	}

	// 创建一个缓冲区来保存模板生成的结果
	var buf bytes.Buffer
	// 使用模板和数据生成输出
	err = tmpl.Execute(&buf, context)

	chatGptPrompt := buf.String()
	return chatGptPrompt, err
}

type ChatGptResp struct {
	Title   string
	Content string
	Topic   string
	Tags    string
}
