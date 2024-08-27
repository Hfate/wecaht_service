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
	SettStatus     int    `json:"sett_status"`
	SettledRevenue int    `json:"settled_revenue"`
	SettNo         string `json:"sett_no"`
	MailSendCnt    string `json:"mail_send_cnt"`
	SlotRevenue    string `json:"slot_revenue"`
}

func (WechatSettlement) TableName() string {
	return "wechat_settlement"
}
