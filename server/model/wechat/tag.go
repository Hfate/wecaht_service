package wechat

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type Tag struct {
	global.GVA_MODEL
	Tag string `json:"tag"`
}

func (Tag) TableName() string {
	return "wechat_tag"
}
