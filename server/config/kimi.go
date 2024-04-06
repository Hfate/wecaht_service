package config

type Kimi struct {
	ApiUrl       string `mapstructure:"api-url" json:"api-url" yaml:"api-url"`
	RefreshToken string `mapstructure:"refresh-token" json:"refresh-token" yaml:"refresh-token"`
}
