package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/wechat"

type ArticleResponse struct {
	Article wechat.Article `json:"article"`
}
