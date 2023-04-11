// 自动生成模板Stat
package v2ray

// Stat 结构体
type Stat struct {
	ID        uint   `gorm:"primarykey"` // 主键ID
	Category  string `json:"category" form:"category" gorm:"column:category;index:idx_cate_tag;index:idx_cate_ip;comment:user inbound outbound;"`
	Tag       string `json:"tag" form:"tag" gorm:"column:tag;uniqueIndex:udx_tag_time_ip;index:idx_cate_tag;comment:如果category是user，tag就是uid;"`
	Down      uint64 `json:"down" form:"down" gorm:"column:down;comment:字节;"`
	Up        uint64 `json:"up" form:"up" gorm:"column:up;comment:;"`
	ServerIp  string `json:"server_ip" form:"server_ip" gorm:"column:server_ip;uniqueIndex:udx_tag_time_ip;index:idx_cate_ip;"`
	CreatedAt int    `json:"created_at" form:"created_at" gorm:"column:created_at;uniqueIndex:udx_tag_time_ip;index"` // 创建时间
}

// TableName Stat 表名
func (Stat) TableName() string {
	return "v2ray_stat"
}
