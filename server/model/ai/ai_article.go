package ai

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type AIArticle struct {
	global.BASEMODEL
	OriginId          uint64 `json:"originId"`
	Title             string `json:"title"` //
	PortalName        string `json:"portalName"`
	Topic             string `json:"topic"` //
	TargetAccountName string `json:"targetAccountName"`
	TargetAccountId   string `json:"targetAccountId"`
	AuthorName        string `json:"authorName" ` //
	Link              string `json:"link"`        //
	PublishTime       string `json:"publishTime"`
	LikeNum           int    `json:"likeNum"`
	ReadNum           int    `json:"readNum"`
	CommentNum        int    `json:"commentNum"`
	Content           string `json:"content"`
	Tags              string `json:"tags"`
	AuditStatus       int    `json:"auditStatus"`
}

func (AIArticle) TableName() string {
	return "ai_article"
}
