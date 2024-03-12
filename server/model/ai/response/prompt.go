package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/ai"

type PromptResponse struct {
	Prompt ai.Prompt `json:"prompt"`
}
