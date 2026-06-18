package config

type Server struct {
	JWT     JWT     `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Zap     Zap     `mapstructure:"zap" json:"zap" yaml:"zap"`
	Redis   Redis   `mapstructure:"redis" json:"redis" yaml:"redis"`
	Email   Email   `mapstructure:"email" json:"email" yaml:"email"`
	System  System  `mapstructure:"system" json:"system" yaml:"system"`
	Captcha Captcha `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
	// auto
	AutoCode Autocode `mapstructure:"autocode" json:"autocode" yaml:"autocode"`
	// gorm
	Mysql  Mysql           `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Mssql  Mssql           `mapstructure:"mssql" json:"mssql" yaml:"mssql"`
	Pgsql  Pgsql           `mapstructure:"pgsql" json:"pgsql" yaml:"pgsql"`
	Oracle Oracle          `mapstructure:"oracle" json:"oracle" yaml:"oracle"`
	DBList []SpecializedDB `mapstructure:"db-list" json:"db-list" yaml:"db-list"`
	// oss
	Local      Local      `mapstructure:"local" json:"local" yaml:"local"`
	Qiniu      Qiniu      `mapstructure:"qiniu" json:"qiniu" yaml:"qiniu"`
	AliyunOSS  AliyunOSS  `mapstructure:"aliyun-oss" json:"aliyun-oss" yaml:"aliyun-oss"`
	HuaWeiObs  HuaWeiObs  `mapstructure:"hua-wei-obs" json:"hua-wei-obs" yaml:"hua-wei-obs"`
	TencentCOS TencentCOS `mapstructure:"tencent-cos" json:"tencent-cos" yaml:"tencent-cos"`
	AwsS3      AwsS3      `mapstructure:"aws-s3" json:"aws-s3" yaml:"aws-s3"`

	Excel Excel `mapstructure:"excel" json:"excel" yaml:"excel"`
	Timer Timer `mapstructure:"timer" json:"timer" yaml:"timer"`

	// 跨域配置
	Cors CORS `mapstructure:"cors" json:"cors" yaml:"cors"`

	STAT_PORT              int64        `mapstructure:"stat_port" json:"stat_port" yaml:"stat_port"`
	TrafficCollectInterval string       `mapstructure:"traffic_collect_interval" json:"traffic_collect_interval" yaml:"traffic_collect_interval"`
	SysInfoCollectInterval string       `mapstructure:"sysinfo_collect_interval" json:"sysinfo_collect_interval" yaml:"sysinfo_collect_interval"`
	TrafficMeter           TrafficMeter `mapstructure:"traffic-meter" json:"traffic-meter" yaml:"traffic-meter"`
	SiliconFlow            SiliconFlow  `mapstructure:"silicon-flow" json:"silicon-flow" yaml:"silicon-flow"`

	// BWG VPS 管理配置
	BWG BWG `mapstructure:"bwg" json:"bwg" yaml:"bwg"`
}

// BWG VPS 管理配置
type BWG struct {
	VeID   string `mapstructure:"veid" json:"veid" yaml:"veid"`
	ApiKey string `mapstructure:"apiKey" json:"apiKey" yaml:"apiKey"`
}
