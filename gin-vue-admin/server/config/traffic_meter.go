package config

type TrafficMeter struct {
	Enable        bool   `mapstructure:"enable" json:"enable" yaml:"enable"`
	StatURL       string `mapstructure:"stat-url" json:"stat-url" yaml:"stat-url"`
	Tag           string `mapstructure:"tag" json:"tag" yaml:"tag"`
	FlushInterval string `mapstructure:"flush-interval" json:"flush-interval" yaml:"flush-interval"`
}
