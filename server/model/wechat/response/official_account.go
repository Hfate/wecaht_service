package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/wechat"

type OfficialAccountResponse struct {
	OfficialAccount wechat.OfficialAccount `json:"officialAccount"`
}
