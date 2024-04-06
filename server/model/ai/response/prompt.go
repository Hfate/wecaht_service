package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type PromptResp struct {
	global.BASEMODEL
	Topic      string   `json:"topic"`
	PromptType int      `json:"promptType"`
	PromptList []string `json:"promptList"`
	Language   string   `json:"language"`
}
