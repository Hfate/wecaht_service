package ai

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type Template struct {
	global.BASEMODEL
	AccountName   string `json:"accountName"`
	AccountId     string `json:"accountId"`
	TemplateValue string `json:"templateValue"`
}

func (Template) TableName() string {
	return "wechat_template"
}
