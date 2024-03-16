package ai

import (
	"bytes"
	"context"
	"github.com/baidubce/bce-qianfan-sdk/go/qianfan"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"go.uber.org/zap"
	"strings"
	"text/template"
)

type QianfanService struct {
}

var QianfanServiceApp = new(QianfanService)

func init() {
	qianfan.GetConfig().AccessKey = "ALTAK5HSinZtO6tas6f0l7und9"
	qianfan.GetConfig().SecretKey = "d4d47be09aef4ff4bbe84564c37bfaa9"
}

func (*QianfanService) Recreation(article ai.Article) ChatGptResp {
	// 可以通过 WithModel 指定模型
	chat := qianfan.NewChatCompletion()

	chatGptPrompt := QianfanServiceApp.parseContentPrompt(article)

	resp, err := chat.Do(
		context.TODO(),
		&qianfan.ChatCompletionRequest{
			Messages: []qianfan.ChatCompletionMessage{
				qianfan.ChatCompletionUserMessage(chatGptPrompt),
			},
		},
	)

	if err != nil {
		global.GVA_LOG.Info("chat gpt响应失败", zap.Error(err))
		return ChatGptResp{}
	}

	content := resp.Result
	content = strings.ReplaceAll(content, "```json", "")
	content = strings.ReplaceAll(content, "```", "")
	content = strings.ReplaceAll(content, "\n", "")
	content = strings.ReplaceAll(content, "”", "\"")
	chatGptResp := ChatGptResp{}

	err = utils.JsonStrToStruct(content, &chatGptResp)
	if err != nil {
		global.GVA_LOG.Info("chat gpt响应失败", zap.Error(err), zap.String("content", content))
		return ChatGptResp{}
	}
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

	return chatGptResp
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
