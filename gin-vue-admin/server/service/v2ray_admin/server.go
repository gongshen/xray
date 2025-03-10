package v2ray_admin

import (
	"encoding/json"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/v2ray"
	v2rayReq "github.com/flipped-aurora/gin-vue-admin/server/model/v2ray/request"
)

type ServerService struct {
}

// CreateServer 创建Server记录
// Author [piexlmax](https://github.com/piexlmax)
func (serverService *ServerService) CreateServer(server *v2ray.Server) (err error) {
	// 初始化配置文件
	config := v2ray.NewXRayConfig(server.Port)
	server.Config, err = json.Marshal(config)
	if err != nil {
		return
	}
	err = global.GVA_DB.Create(server).Error
	return err
}

// DeleteServer 删除Server记录
// Author [piexlmax](https://github.com/piexlmax)
func (serverService *ServerService) DeleteServer(server v2ray.Server) (err error) {
	err = global.GVA_DB.Delete(&server).Error
	return err
}

// DeleteServerByIds 批量删除Server记录
// Author [piexlmax](https://github.com/piexlmax)
func (serverService *ServerService) DeleteServerByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]v2ray.Server{}, "id in ?", ids.Ids).Error
	return err
}

// UpdateServer 更新Server记录
// Author [piexlmax](https://github.com/piexlmax)
func (serverService *ServerService) UpdateServer(server v2ray.Server) (err error) {
	err = global.GVA_DB.Save(&server).Error
	return err
}

// GetServer 根据id获取Server记录
// Author [piexlmax](https://github.com/piexlmax)
func (serverService *ServerService) GetServer(id uint) (server v2ray.Server, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&server).Error
	return
}

// GetServerInfoList 分页获取Server记录
// Author [piexlmax](https://github.com/piexlmax)
func (serverService *ServerService) GetServerInfoList(info v2rayReq.ServerSearch) (list []v2ray.Server, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&v2ray.Server{})
	var servers []v2ray.Server
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Ip != "" {
		db = db.Where("ip = ?", info.Ip)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&servers).Error
	return servers, total, err
}

func (serverService *ServerService) GetAllServer() (srvList []*v2ray.Server, err error) {
	err = global.GVA_DB.Find(&srvList).Error
	return
}

func (ServerService *ServerService) UpdateServerUsedQuota(id uint, used uint64) error {
	return global.GVA_DB.Exec("UPDATE v2ray_server set used_quota = used_quota + ? where ID = ?", used, id).Error
}

func (ServerService *ServerService) ResetServerUsedQuota(id uint) error {
	return global.GVA_DB.Exec("UPDATE v2ray_server set used_quota = 0 where ID = ?", id).Error
}

func (serverService *ServerService) UpdateServerConfig(id uint, config json.RawMessage) (err error) {
	err = global.GVA_DB.Table("v2ray_server").Where("ID = ?", id).Update("config", config).Error
	return err
}

func (serverService *ServerService) SaveServerUsedQuotaLog(log *v2ray.ServerQuotaLog) (err error) {
	err = global.GVA_DB.Create(log).Error
	return err
}