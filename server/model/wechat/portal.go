package wechat

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type Portal struct {
	global.GVA_MODEL
	PortalName string `json:"portalName"`
	PortalKey  string `json:"portalKey"`
	ArticleKey string `json:"articleKey"`
	Link       string `json:"link"`
	Topic      string `json:"topic"`
	GraphQuery string `json:"graphQuery"`
	TargetNum  int    `json:"targetNum"`
	Remark     string `json:"remark"`
}

func (Portal) TableName() string {
	return "wechat_portal"
}
