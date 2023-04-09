// 自动生成模板Binding
package v2ray

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

// Binding 结构体
type Binding struct {
	global.GVA_MODEL
	ServerId int            `json:"server_id" form:"server_id" gorm:"column:server_id;comment:服务器id"` // 服务器id
	Server   Server         `json:"server"`
	UserID   int            `json:"user_id" form:"user_id" gorm:"column:user_id;comment:用户id"` // 用户id
	User     system.SysUser `json:"user"`
}

// TableName Binding 表名
func (Binding) TableName() string {
	return "v2ray_binding"
}
