package ai

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type Url struct {
	global.BASEMODEL
	Url string `json:"url"`
}

func (Url) TableName() string {
	return "wechat_url"
}
