package ai

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type Article struct {
	global.GVA_MODEL
	Title       string `json:"title"` //
	PortalName  string `json:"portalName"`
	Topic       string `json:"topic"`       //
	AuthorName  string `json:"authorName" ` //
	Link        string `json:"link"`        //
	PublishTime string `json:"publishTime"`
	LikeNum     int    `json:"likeNum"`
	ReadNum     int    `json:"readNum"`
	CommentNum  int    `json:"commentNum"`
	Content     string `json:"content"`
	Tags        string `json:"tags"`
}

func (Article) TableName() string {
	return "wechat_article"
}
