package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/wechat"
)

type PortalResponse struct {
	Portal wechat.Portal `json:"portal"`
}
