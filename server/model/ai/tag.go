package ai

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type Tag struct {
	global.BASEMODEL
	Tag string `json:"tag"`
}

func (Tag) TableName() string {
	return "wechat_tag"
}
