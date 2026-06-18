package config

type SiliconFlow struct {
	APIKey  string `mapstructure:"api-key" json:"api-key" yaml:"api-key"`
	BaseURL string `mapstructure:"base-url" json:"base-url" yaml:"base-url"`
	Model   string `mapstructure:"model" json:"model" yaml:"model"`
	Timeout string `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
}
