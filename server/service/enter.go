package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
)

type ServiceGroup struct {
	SystemServiceGroup system.ServiceGroup
	AIServiceGroup     ai.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
