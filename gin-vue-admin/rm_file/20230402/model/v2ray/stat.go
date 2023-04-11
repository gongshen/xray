// 自动生成模板Stat
package v2ray

// Stat 结构体
type Stat struct {
	ID         uint   `gorm:"primarykey"` // 主键ID
	Category   string `json:"category" form:"category" gorm:"column:category;comment:inbound、outbound、user;"`
	Tag        string `json:"tag" form:"tag" gorm:"column:tag;uniqueIndex:udx_tag_day;comment:如果category是user，tag就是uid;"`
	Down       int    `json:"down" form:"down" gorm:"column:down;comment:字节;"`
	Up         int    `json:"up" form:"up" gorm:"column:up;comment:;"`
	CreatedDay int    `json:"created_day" form:"created_day" gorm:"column:created_day;uniqueIndex:udx_tag_day;comment:;"`
}

// TableName Stat 表名
func (Stat) TableName() string {
	return "user_stat"
}
