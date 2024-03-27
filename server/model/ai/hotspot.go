package ai

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type Hotspot struct {
	global.BASEMODEL
	Headline   string `json:"headline"`
	PortalName string `json:"portalName"`
	Link       string `json:"link"`
	Trending   int    `json:"trending"`
	Topic      string `json:"topic"`
	UseTimes   int    `json:"useTimes"`
}

func (Hotspot) TableName() string {
	return "hotspot"
}
