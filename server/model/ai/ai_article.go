package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

type AIArticle struct {
	global.BASEMODEL
	OriginId          uint64    `json:"originId"`
	Title             string    `json:"title"` //
	PortalName        string    `json:"portalName"`
	Topic             string    `json:"topic"` //
	TargetAccountName string    `json:"targetAccountName"`
	TargetAccountId   string    `json:"targetAccountId"`
	AuthorName        string    `json:"authorName" ` //
	Link              string    `json:"link"`        //
	PublishTime       time.Time `json:"publishTime"`
	LikeNum           int       `json:"likeNum"`
	ReadNum           int       `json:"readNum"`
	CommentNum        int       `json:"commentNum"`
	Content           string    `json:"content"`
	Tags              string    `json:"tags"`
	ArticleStatus     int       `json:"articleStatus"`
	MediaId           string    `json:"mediaId"`
	PublishId         int64     `json:"publishId"`
}

func (AIArticle) TableName() string {
	return "ai_article"
}
