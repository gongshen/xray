// 自动生成模板Server
package v2ray

import (
	"encoding/json"
	"time"
)

// Server 结构体
type Server struct {
	ID        uint            `gorm:"primarykey"` // 主键ID
	Ip        string          `json:"ip" form:"ip" gorm:"column:ip;uniqueIndex:udx_ip;"`
	Remark    string          `json:"remark" form:"remark" gorm:"column:remark;comment:;"`
	Port      int64           `json:"port" form:"port" gorm:"column:port"`
	ResetDate time.Time       `json:"reset_date" form:"reset_date" gorm:"column:reset_date"` // 流量重置日期
	Config    json.RawMessage `json:"config" form:"config" gorm:"column:config;type:json"`
	CreatedAt time.Time       // 创建时间
}

// TableName Server 表名
func (Server) TableName() string {
	return "v2ray_server"
}
