package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/spf13/cast"
	"time"
)

func BaseModel() global.BASEMODEL {
	return global.BASEMODEL{ID: cast.ToString(utils.GenID()), CreatedAt: time.Now(), UpdatedAt: time.Now()}
}
