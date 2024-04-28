package ai

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type OfficialAccount struct {
	global.BASEMODEL
	AccountName       string `json:"accountName"`
	AccountId         string `json:"accountId"`
	Topic             string `json:"topic"`
	UserEmail         string `json:"userEmail"`
	AppId             string `json:"appId"`
	HeadImgUrl        string `json:"headImgUrl"`
	Signature         string `json:"signature"`
	AppSecret         string `json:"appSecret"`
	Token             string `json:"token"`
	EncodingAESKey    string `json:"encodingAesKey"`
	DefaultAuthorName string `json:"defaultAuthorName"`
	Remark            string `json:"remark"`
	TargetNum         int    `json:"targetNum"`
	CreateTypes       string `json:"createTypes"`
	CssFormat         string `json:"cssFormat"`
	NeedOpenComment   uint   `json:"needOpenComment"`
}

func (OfficialAccount) TableName() string {
	return "wechat_official_account"
}
