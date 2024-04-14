package config

type QianFan struct {
	AccessKey string `mapstructure:"access-key" json:"access-key" yaml:"access-key"`
	SecretKey string `mapstructure:"secret-key" json:"secret-key" yaml:"secret-key"`
	Cookie    string `mapstructure:"cookie" json:"cookie" yaml:"cookie"`
}
