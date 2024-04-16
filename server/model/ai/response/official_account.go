package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/ai"

type OfficialAccountResponse struct {
	*ai.OfficialAccount
	CreateTypeList []int `json:"createTypeList"`
}
