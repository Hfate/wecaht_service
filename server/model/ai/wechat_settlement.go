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

type WechatSettlementExcel struct {
	AccountName    string `json:"accountName" xlsx:"Account Name"`
	Date           string `json:"date" xlsx:"Date"`
	Zone           string `json:"zone" xlsx:"Zone"`
	Month          string `json:"month" xlsx:"Month"`
	Order          int    `json:"order" xlsx:"Order"`
	SettStatus     string `json:"settStatus" xlsx:"Sett Status"`
	SettledRevenue string `json:"settledRevenue" xlsx:"Settled Revenue"`
	SettNo         string `json:"settNo" xlsx:"Sett No"`
	MailSendCnt    string `json:"mailSendCnt" xlsx:"Mail Send Cnt"`
	SlotRevenue    string `json:"slotRevenue" xlsx:"Slot Revenue"`
}
