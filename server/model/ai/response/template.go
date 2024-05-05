package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/ai"

type TemplateResponse struct {
	Template ai.Template `json:"template"`
}
