package ai

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type Subject struct {
	global.BASEMODEL
	Topic    string `json:"topic"`
	Subject  string `json:"subject"`
	UseTimes int    `json:"useTimes"`
}

func (Subject) TableName() string {
	return "wechat_subject"
}
