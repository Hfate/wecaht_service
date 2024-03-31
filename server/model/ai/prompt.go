package ai

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type Prompt struct {
	global.BASEMODEL
	Topic      string `json:"topic"`
	PromptType int    `json:"promptType"`
	Prompt     string `json:"prompt"`
	Language   string `json:"language"`
}

func (Prompt) TableName() string {
	return "wechat_prompt"
}

// 定义一个名为Color的“枚举”
const (
	ContentRecreation = 1 // iota从0开始，每定义一个常量自动加1
	TitleRecreation   = 2
	HotSpotWrite      = 3
	TopicWrite        = 4
	TitleCreate       = 5
	AddImage          = 6
)
