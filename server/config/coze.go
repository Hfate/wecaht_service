package config

type Coze struct {
	ApiUrl      string `mapstructure:"api-url" json:"api-url" yaml:"api-url"`
	AccessToken string `mapstructure:"access-token" json:"access-token" yaml:"access-token"`
	BotId       string `mapstructure:"bot-id" json:"bot-id" yaml:"bot-id"`
}
