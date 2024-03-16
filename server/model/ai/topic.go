package ai

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type Topic struct {
	global.BASEMODEL
	Topic string `json:"topic"`
}

func (Topic) TableName() string {
	return "wechat_topic"
}
