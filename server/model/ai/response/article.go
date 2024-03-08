package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/ai"

type ArticleResponse struct {
	Article ai.Article `json:"article"`
}
