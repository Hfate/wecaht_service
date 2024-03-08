package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type JwtBlacklist struct {
	global.BASEMODEL
	Jwt string `gorm:"type:text;comment:jwt"`
}
