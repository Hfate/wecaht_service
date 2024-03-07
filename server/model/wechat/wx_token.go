package wechat

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type WxToken struct {
	global.GVA_MODEL
	SlaveUser  string `json:"slaveUser"`
	SlaveSid   string `json:"slaveSid"`
	BizUin     string `json:"bizUin"`
	DataTicket string `json:"dataTicket"`
	RandInfo   string `json:"randInfo"`
	PassTicket string `json:"passTicket"`
}

func (WxToken) TableName() string {
	return "wechat_token"
}
