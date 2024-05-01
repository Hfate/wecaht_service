package ai

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type Article struct {
	global.BASEMODEL
	Title       string `json:"title"` //
	PortalName  string `json:"portalName"`
	Topic       string `json:"topic" xlsx:"Topic"` //
	AuthorName  string `json:"authorName" `        //
	Link        string `json:"link" xlsx:"Link"`   //
	PublishTime string `json:"publishTime"`
	LikeNum     int    `json:"likeNum"`
	ReadNum     int    `json:"readNum"`
	CommentNum  int    `json:"commentNum"`
	Content     string `json:"content"`
	Comment     string `json:"comment"`
	Tags        string `json:"tags"`
	UseTimes    int    `json:"useTimes"`
	HotspotId   uint64 `json:"hotspotId"`
	IsHot       bool   `json:"isHot"`
}

func (Article) TableName() string {
	return "wechat_article"
}

type ArticleExcl struct {
	Title       string `json:"title"  xlsx:"Title"` //
	PortalName  string `json:"portalName" xlsx:"PortalName"`
	Topic       string `json:"topic" xlsx:"Topic"`           //
	AuthorName  string `json:"authorName" xlsx:"AuthorName"` //
	Link        string `json:"link" xlsx:"Link"`             //
	PublishTime string `json:"publishTime" xlsx:"PublishTime"`
	LikeNum     int    `json:"likeNum" xlsx:"LikeNum"`
	ReadNum     int    `json:"readNum" xlsx:"ReadNum"`
	CommentNum  int    `json:"commentNum" xlsx:"CommentNum"`
	Content     string `json:"content" xlsx:"Content"`
	Tags        string `json:"tags" xlsx:"Tags"`
}

type ArticleExclUpload struct {
	Title       string `json:"title"  xlsx:"Title"`
	Topic       string `json:"topic" xlsx:"Topic"`
	Link        string `json:"link" xlsx:"Link"`
	ReadNum     string `json:"readNum" xlsx:"ReadNum"`
	Comment     string `json:"comment" xlsx:"Comment"`
	Content     string `json:"content" xlsx:"Content"`
	LikeNum     string `json:"LikeNum" xlsx:"LikeNum"`
	PortalName  string `json:"portalName" xlsx:"PortalName"`
	PublishTime string `json:"publishTime" xlsx:"PublishTime"`
}

type ArticleStats struct {
	Topic string `json:"topic"`
	Count int64  `json:"count"`
}
