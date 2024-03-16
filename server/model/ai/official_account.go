package ai

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type OfficialAccount struct {
	global.BASEMODEL
	AccountName string `json:"accountName"`
	AccountId   string `json:"accountId"`
	Topic       string `json:"topic"`
	UserEmail   string `json:"userEmail"`
	AppId       string `json:"appId"`
	AppSecret   string `json:"appSecret"`
	Token       string `json:"token"`
	Remark      string `json:"remark"`
}

func (OfficialAccount) TableName() string {
	return "wechat_official_account"
}
