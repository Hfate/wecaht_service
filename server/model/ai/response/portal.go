package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/wechat"
)

type PortalResponse struct {
	Portal ai.Portal `json:"portal"`
}
