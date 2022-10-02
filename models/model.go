package models

type Traffic struct {
	Id     int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Tag    string `json:"tag" form:"tag" gorm:"unique"`
	Limit  int64  `json:"limit"`                             // 流量限制
	Used   int64  `json:"used"`                              // 使用流量
	Enable bool   `json:"enable" form:"enable" gorm:"index"` // 是否允许使用
}
