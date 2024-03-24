package ai

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type Media struct {
	global.BASEMODEL
	Topic             string `json:"topic"`
	MediaID           string `json:"mediaID"`
	Link              string `json:"link"`
	FileName          string `json:"fileName"`
	TargetAccountName string `json:"targetAccountName"`
	TargetAccountId   string `json:"targetAccountId"`
	SeqNum            int    `json:"seqNum"`
	Tag               string `json:"tag"`
}

func (Media) TableName() string {
	return "wechat_media"
}
