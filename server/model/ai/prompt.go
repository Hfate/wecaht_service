package ai

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type Prompt struct {
	global.BASEMODEL
	Topic      string `json:"topic"`
	PromptType string `json:"promptType"`
	Prompt     string `json:"prompt"`
	Language   string `json:"language"`
}

func (Prompt) TableName() string {
	return "wechat_prompt"
}
