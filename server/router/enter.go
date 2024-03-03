package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/router/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/router/example"
	"github.com/flipped-aurora/gin-vue-admin/server/router/system"
)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
	AI      ai.AIGroup
}

var RouterGroupApp = new(RouterGroup)
