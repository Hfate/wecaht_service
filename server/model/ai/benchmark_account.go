package ai

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type BenchmarkAccount struct {
	global.BASEMODEL
	AccountName string `json:"accountName"`
	AccountId   string `json:"accountId"`
	Topic       string `json:"topic"`
	InitNum     int    `json:"initNum"`
	ArticleLink string `json:"articleLink"`
	Key         string `json:"key"`
}

func (BenchmarkAccount) TableName() string {
	return "wechat_benchmark_account"
}
