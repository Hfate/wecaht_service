package config

type ChatModel struct {
	ApiUrl       string `mapstructure:"api-url" json:"api-url" yaml:"api-url"`
	RefreshToken string `mapstructure:"refresh-token" json:"refresh-token" yaml:"refresh-token"`
	Model        string `mapstructure:"model" json:"model" yaml:"model"`
	ModelType    string `mapstructure:"model-type" json:"model-type" yaml:"model-type"`
}
