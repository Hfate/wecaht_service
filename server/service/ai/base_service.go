package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"time"
)

func BaseModel() global.BASEMODEL {
	return global.BASEMODEL{ID: utils.GenID(), CreatedAt: time.Now(), UpdatedAt: time.Now()}
}
