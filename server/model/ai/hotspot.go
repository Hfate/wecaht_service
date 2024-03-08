package ai

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type Hotspot struct {
	global.GVA_MODEL
	Headlines  string `json:"headlines"`
	PortalName string `json:"portalName"`
	Link       string `json:"link"`
}
