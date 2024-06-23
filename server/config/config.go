package config

type Server struct {
	JWT        JWT         `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Zap        Zap         `mapstructure:"zap" json:"zap" yaml:"zap"`
	Redis      Redis       `mapstructure:"redis" json:"redis" yaml:"redis"`
	Sitong     Sitong      `mapstructure:"sitong" json:"sitong" yaml:"sitong"`
	Kimi       Kimi        `mapstructure:"kimi" json:"kimi" yaml:"kimi"`
	Dajiala    Dajiala     `mapstructure:"dajiala" json:"dajiala" yaml:"dajiala"`
	ChatModels []ChatModel `mapstructure:"chat-models" json:"chat-models" yaml:"chat-models"`
	Baidu      Baidu       `mapstructure:"baidu" json:"baidu" yaml:"baidu"`
	Mongo      Mongo       `mapstructure:"mongo" json:"mongo" yaml:"mongo"`
	Email      Email       `mapstructure:"email" json:"email" yaml:"email"`
	System     System      `mapstructure:"system" json:"system" yaml:"system"`
	Captcha    Captcha     `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
	// auto
	AutoCode Autocode `mapstructure:"autocode" json:"autocode" yaml:"autocode"`
	// gorm
	Mysql  Mysql           `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Mssql  Mssql           `mapstructure:"mssql" json:"mssql" yaml:"mssql"`
	Pgsql  Pgsql           `mapstructure:"pgsql" json:"pgsql" yaml:"pgsql"`
	Oracle Oracle          `mapstructure:"oracle" json:"oracle" yaml:"oracle"`
	Sqlite Sqlite          `mapstructure:"sqlite" json:"sqlite" yaml:"sqlite"`
	DBList []SpecializedDB `mapstructure:"db-list" json:"db-list" yaml:"db-list"`
	// oss
	Local      Local      `mapstructure:"local" json:"local" yaml:"local"`
	Qiniu      Qiniu      `mapstructure:"qiniu" json:"qiniu" yaml:"qiniu"`
	AliyunOSS  AliyunOSS  `mapstructure:"aliyun-oss" json:"aliyun-oss" yaml:"aliyun-oss"`
	HuaWeiObs  HuaWeiObs  `mapstructure:"hua-wei-obs" json:"hua-wei-obs" yaml:"hua-wei-obs"`
	TencentCOS TencentCOS `mapstructure:"tencent-cos" json:"tencent-cos" yaml:"tencent-cos"`
	AwsS3      AwsS3      `mapstructure:"aws-s3" json:"aws-s3" yaml:"aws-s3"`

	Excel Excel `mapstructure:"excel" json:"excel" yaml:"excel"`
	Coze  Coze  `mapstructure:"coze" json:"coze" yaml:"coze"`

	// 跨域配置
	Cors CORS `mapstructure:"cors" json:"cors" yaml:"cors"`
}
