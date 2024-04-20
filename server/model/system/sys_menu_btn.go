package system

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type SysBaseMenuBtn struct {
	global.BASEMODEL
	Name          string `json:"name" gorm:"comment:按钮关键key"`
	Desc          string `json:"desc" gorm:"按钮备注"`
	SysBaseMenuID string `json:"sysBaseMenuID" gorm:"comment:菜单ID"`
}
