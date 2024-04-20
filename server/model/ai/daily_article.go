package ai

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type DailyArticle struct {
	global.BASEMODEL
	Title             string `json:"title"` //
	PortalName        string `json:"portalName"`
	Topic             string `json:"topic" xlsx:"Topic"` //
	AuthorName        string `json:"authorName" `        //
	Link              string `json:"link" xlsx:"Link"`   //
	PublishTime       string `json:"publishTime"`
	LikeNum           int    `json:"likeNum"`
	ReadNum           int    `json:"readNum"`
	CommentNum        int    `json:"commentNum"`
	Comment           string `json:"comment"`
	Content           string `json:"content"`
	Tags              string `json:"tags"`
	UseTimes          int    `json:"useTimes"`
	HotspotId         uint64 `json:"hotspotId"`
	BatchId           string `json:"batchId"`
	TargetAccountId   string `json:"targetAccountId"`
	TargetAccountName string `json:"targetAccountName"`
	IsHot             bool   `json:"isHot"`
}

func (DailyArticle) TableName() string {
	return "wechat_daily_article"
}
