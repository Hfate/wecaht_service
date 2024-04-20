package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/ai"

type DailyArticleResponse struct {
	DailyArticle ai.DailyArticle `json:"article"`
}

type DailyArticleParentResponse struct {
	ai.DailyArticle
	Children []ai.DailyArticle `json:"children"`
}
