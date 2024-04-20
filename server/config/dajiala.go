package config

type Dajiala struct {
	Key              string `mapstructure:"key" json:"key" yaml:"key"`
	GetRankUrl       string `mapstructure:"get-rank-url" json:"get-rank-url" yaml:"get-rank-url"`
	PostHistoryUrl   string `mapstructure:"post-history-url" json:"post-condition-url" yaml:"post-condition-url"`
	ArticleDetailUrl string `mapstructure:"article-detail-url" json:"article-detail-url" yaml:"article-detail-url"`
	ReadAndZanUrl    string `mapstructure:"read-and-zan-url" json:"read-and-zan-url" yaml:"read-and-zan-url"`
}
