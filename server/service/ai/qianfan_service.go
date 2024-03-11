package ai

import (
	"context"
	"fmt"
	"github.com/baidubce/bce-qianfan-sdk/go/qianfan"
	"strings"
)

type QianfanService struct {
}

var QianfanServiceApp = new(QianfanService)

func init() {
	qianfan.GetConfig().AccessKey = "ALTAK5HSinZtO6tas6f0l7und9"
	qianfan.GetConfig().SecretKey = "d4d47be09aef4ff4bbe84564c37bfaa9"
}

func (*QianfanService) Recreation(title, content string) (string, string) {
	title = QianfanServiceApp.TitleRecreation(title)
	content = QianfanServiceApp.ContentRecreation(content)
	return title, content
}

func (*QianfanService) TitleRecreation(title string) string {
	// 可以通过 WithModel 指定模型
	chat := qianfan.NewChatCompletion()

	resp, err := chat.Do(
		context.TODO(),
		&qianfan.ChatCompletionRequest{
			Messages: []qianfan.ChatCompletionMessage{
				qianfan.ChatCompletionUserMessage("你是一个微信公众号的资深流量主，热爱生活，热衷于向你的粉丝投递好的有深度的内容。" +
					"请基于以下{}提供的视频原始标题进行改写为高质量标题，在保证原始视频标题不丢失的前提下，改写的目的是尽可能的让用户看到标题后就能够点击进来浏览视频以提升视频的点击率" +
					"\n请按照如下要求对视频标题进行改写：" +
					"\n1.必须基于{}中的原始标题进行改写，不能缺少原始标题内容" +
					"\n2.也不能生成与原有标题不相关的内容" +
					"\n3.原始视频标题与改写后的高质量标题必须一一对应" +
					"\n4.生成的标题尽可能的吸引用户但又不能过于标题党" +
					"\n5.每次生成1个标题，每个标题不超过20字" +
					"\n### 示例 ###\n1. 原始视频标题：{108集唱歌教程} 改写后标题：108集唱歌教程，让你秒变唱将\n2. 原始视频标题：{108集唱歌教程} 改写后标题：108集唱歌教程，教你如何快速提升唱功\n3. 原始视频标题：{灌篮高手，全国大赛} 改写后标题：灌篮高手，全国大赛，湘北队何去何从\n4. 原始视频标题：{灌篮高手，全国大赛} 改写后标题：灌篮高手，全国大赛，樱木最后的灌篮\n\n#### 视频原始标题 ####\n视频原始标题：" + title + "\n\n改写后标题："),
			},
		},
	)
	if err != nil {
		fmt.Print(err)
	}
	result := resp.Result
	result = strings.ReplaceAll(result, "{}", "")

	return resp.Result
}

func (*QianfanService) ContentRecreation(content string) string {

	// 可以通过 WithModel 指定模型
	chat := qianfan.NewChatCompletion()

	resp, err := chat.Do(
		context.TODO(),
		&qianfan.ChatCompletionRequest{
			Messages: []qianfan.ChatCompletionMessage{
				qianfan.ChatCompletionUserMessage("用1种不同的方式改写以下段落，以避免重复，同时保持其含义：" + content),
			},
		},
	)
	if err != nil {
		fmt.Print(err)
	}

	return resp.Result
}
