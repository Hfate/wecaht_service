package config

type Coze struct {
	ApiUrl          string `mapstructure:"api-url" json:"api-url" yaml:"api-url"`
	AccessToken     string `mapstructure:"access-token" json:"access-token" yaml:"access-token"`
	HotBotId        string `mapstructure:"hot-bot-id" json:"hot-bot-id" yaml:"hot-bot-id"`
	RecreationBotId string `mapstructure:"recreation-bot-id" json:"recreation-bot-id" yaml:"recreation-bot-id"`
}
