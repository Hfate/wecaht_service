package wechat

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type BenchmarkAccount struct {
	global.GVA_MODEL
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
