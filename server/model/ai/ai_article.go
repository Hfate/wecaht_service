package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

type AIArticle struct {
	global.BASEMODEL
	BatchId           string    `json:"batchId"`
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
	ArticleStatus     int       `json:"articleStatus"` // 0  刚生成  1 发送至草稿箱  2 发布成功 3 群发成功  4 发布失败
	ErrMessage        string    `json:"errMessage"`
	MediaId           string    `json:"mediaId"`
	PublishId         int64     `json:"publishId"`
	MsgDataID         int64     `json:"msgDataId"`
	MsgId             int64     `json:"msgId"`
	Params            string    `json:"params"`
	ProcessStatus     int       `json:"processStatus"`
	ProcessParams     string    `json:"processParams"`
	Percent           int       `json:"percent"`
}

func (AIArticle) TableName() string {
	return "ai_article"
}

// 定义一个名为Color的“枚举”
const (
	ProcessInit      = 0 // iota从0开始，每定义一个常量自动加1
	ProcessCreateIng = 10
	ProcessAddImgIng = 20
	ProcessCreated   = 30
)
