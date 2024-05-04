package config

type Sitong struct {
	AccessKey string `mapstructure:"access-key" json:"access-key" yaml:"access-key"`
	SecretKey string `mapstructure:"secret-key" json:"secret-key" yaml:"secret-key"`
	ApiUrl    string `mapstructure:"api-url" json:"api-url" yaml:"api-url"`
}
