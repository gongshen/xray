// 自动生成模板Binding
package v2ray

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"time"
)

// Binding 结构体
type Binding struct {
	ID        uint           `gorm:"primarykey"` // 主键ID
	CreatedAt time.Time      // 创建时间
	UserID    int            `json:"user_id" form:"user_id" gorm:"column:user_id;uniqueIndex:udx_uid_sid;comment:用户id"` // 用户id
	User      system.SysUser `json:"user"`
	ServerID  int            `json:"server_id" form:"server_id" gorm:"column:server_id;uniqueIndex:udx_uid_sid;index"`
	Server    Server         `json:"server"`
	AlterID   int64          `json:"alter_id" form:"alter_id" gorm:"column:alter_id"`
	Level     int64          `json:"level" form:"level" gorm:"column:level;default:0"`
}

// TableName Binding 表名
func (Binding) TableName() string {
	return "v2ray_binding"
}
