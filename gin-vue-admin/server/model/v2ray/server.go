// 自动生成模板Server
package v2ray

import (
	"encoding/json"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Server 结构体
type Server struct {
	ID         uint            `json:"ID" gorm:"primarykey"` // 主键ID
	Ip         string          `json:"ip" form:"ip" gorm:"column:ip;uniqueIndex:udx_ip;"`
	Remark     string          `json:"remark" form:"remark" gorm:"column:remark;comment:;"`
	Port       int64           `json:"port" form:"port" gorm:"column:port"`
	StatPort   int             `json:"stat_port" form:"stat_port" gorm:"column:stat_port;default:0"` // 统计服务端口，0 表示使用全局配置
	ResetDate  int             `json:"reset_date" form:"reset_date" gorm:"column:reset_date"`        // 流量重置日期
	Config     json.RawMessage `json:"config" form:"config" gorm:"column:config;type:json"`
	CreatedAt  time.Time       // 创建时间
	UsedQuota  uint64          `json:"used_quota" form:"used_quota" gorm:"column:used_quota"`    // 服务器已使用额度
	TotalQuota uint64          `json:"total_quota" form:"total_quota" gorm:"column:total_quota"` // 总额度
	// 系统信息 (由 stat 程序上报)
	DiskTotal  uint64  `json:"disk_total" gorm:"column:disk_total"`   // 磁盘总量 (MB)
	DiskUsed   uint64  `json:"disk_used" gorm:"column:disk_used"`     // 磁盘已用 (MB)
	MemTotal   uint64  `json:"mem_total" gorm:"column:mem_total"`     // 内存总量 (MB)
	MemUsed    uint64  `json:"mem_used" gorm:"column:mem_used"`       // 内存已用 (MB)
	CPUPercent float64 `json:"cpu_percent" gorm:"column:cpu_percent"` // CPU使用率
	SysInfoAt  int64   `json:"sysinfo_at" gorm:"column:sysinfo_at"`   // 系统信息更新时间
}

// TableName Server 表名
func (Server) TableName() string {
	return "v2ray_server"
}

// GetStatPort 返回有效的统计服务端口
// 如果 StatPort > 0，返回 StatPort；否则返回全局配置
func (s *Server) GetStatPort() int {
	if s.StatPort > 0 {
		return s.StatPort
	}
	return int(global.GVA_CONFIG.STAT_PORT)
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
