package models

type Traffic struct {
	Id     int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Tag    string `json:"tag" form:"tag" gorm:"unique"`
	Used   int64  `json:"used" form:"used"`                  // 当前使用流量
	Base   int64  `json:"base" form:"base"`                  // 基础流量
	Enable bool   `json:"enable" form:"enable" gorm:"index"` // 是否启用
}

func (*Traffic) TableName() string {
	return "traffics"
}

type TrafficHistory struct {
	Id         int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Used       int64  `json:"used" form:"used"`
	Tag        string `json:"tag" form:"tag" gorm:"uniqueIndex:udx_create_tag;index:idx:idx_tag"`
	CreatedDay string `json:"created_day" form:"created_day"  gorm:"uniqueIndex:udx_create_tag"`
}

func (*TrafficHistory) TableName() string {
	return "traffic_histories"
}
