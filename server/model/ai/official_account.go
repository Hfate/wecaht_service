package ai

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type OfficialAccount struct {
	global.BASEMODEL
	AccountName string `json:"accountName"`
	AccountId   string `json:"accountId"`
	Topic       string `json:"topic"`
	UserEmail   string `json:"userEmail"`
	WxToken     string `json:"wxToken"`
	Remark      string `json:"remark"`
}

func (OfficialAccount) TableName() string {
	return "wechat_official_account"
}
