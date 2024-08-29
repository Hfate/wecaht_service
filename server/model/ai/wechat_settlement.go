package ai

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type WechatSettlement struct {
	global.BASEMODEL
	AccountName    string `json:"accountName"`
	AccountId      string `json:"accountId"`
	Date           string `json:"date"`
	Zone           string `json:"zone"`
	Month          string `json:"month"`
	Order          int    `json:"order"`
	SettStatus     int    `json:"settStatus"`
	SettledRevenue int    `json:"settledRevenue"`
	SettNo         string `json:"settNo"`
	MailSendCnt    string `json:"mailSendCnt"`
	SlotRevenue    string `json:"slotRevenue"`
}

func (WechatSettlement) TableName() string {
	return "wechat_settlement"
}
