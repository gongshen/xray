// 自动生成模板Server
package v2ray

import (
	"encoding/json"
	"time"
)

// Server 结构体
type Server struct {
	ID         uint            `json:"ID" gorm:"primarykey"` // 主键ID
	Ip         string          `json:"ip" form:"ip" gorm:"column:ip;uniqueIndex:udx_ip;"`
	Remark     string          `json:"remark" form:"remark" gorm:"column:remark;comment:;"`
	Port       int64           `json:"port" form:"port" gorm:"column:port"`
	ResetDate  int             `json:"reset_date" form:"reset_date" gorm:"column:reset_date"` // 流量重置日期
	Config     json.RawMessage `json:"config" form:"config" gorm:"column:config;type:json"`
	CreatedAt  time.Time       // 创建时间
	UsedQuota  uint64          `json:"used_quota" form:"used_quota" gorm:"column:used_quota"`    // 服务器已使用额度
	TotalQuota uint64          `json:"total_quota" form:"total_quota" gorm:"column:total_quota"` // 总额度
}

// TableName Server 表名
func (Server) TableName() string {
	return "v2ray_server"
}

// 服务器额度使用历史记录，每月存一条
type ServerQuotaLog struct {
	ID        uint
	ServerID  int    `json:"server_id" form:"server_id" gorm:"column:server_id;uniqueIndex:udx_server_time"`
	Server    Server `json:"server"`
	CreatedAt int64  `json:"created_at" form:"created_at" gorm:"column:created_at;uniqueIndex:udx_server_time;index"` // 创建时间
	UsedQuota uint64 `json:"used_quota" form:"used_quota" gorm:"column:used_quota"`
}

func (ServerQuotaLog) TableName() string {
	return "v2ray_server_used_quota_log"
}
