package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/ai"

type OfficialAccountResponse struct {
	OfficialAccount ai.OfficialAccount `json:"officialAccount"`
}
