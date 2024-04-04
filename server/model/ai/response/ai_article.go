package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
)

type AIArticleResponse struct {
	AIArticle ai.AIArticle `json:"article"`
}

type AIArticleParentResponse struct {
	ai.AIArticle
	Children []ai.AIArticle `json:"children"`
}
