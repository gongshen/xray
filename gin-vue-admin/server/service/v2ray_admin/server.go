package v2ray_admin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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

// SysInfo 系统信息结构 (与 stat 程序返回的结构对应)
type SysInfo struct {
	DiskTotal  uint64  `json:"dt"` // 磁盘总量 (MB)
	DiskUsed   uint64  `json:"du"` // 磁盘已用 (MB)
	MemTotal   uint64  `json:"mt"` // 内存总量 (MB)
	MemUsed    uint64  `json:"mu"` // 内存已用 (MB)
	CPUPercent float64 `json:"cp"` // CPU使用率
	Timestamp  int64   `json:"ts"` // 时间戳
}

// UpdateServerSysInfo 更新服务器系统信息
func (serverService *ServerService) UpdateServerSysInfo(serverID uint, info *SysInfo) error {
	return global.GVA_DB.Table("v2ray_server").Where("ID = ?", serverID).Updates(map[string]interface{}{
		"disk_total":  info.DiskTotal,
		"disk_used":   info.DiskUsed,
		"mem_total":   info.MemTotal,
		"mem_used":    info.MemUsed,
		"cpu_percent": info.CPUPercent,
		"sysinfo_at":  info.Timestamp,
	}).Error
}

// RestartVPS 重启 VPS 服务器
func (serverService *ServerService) RestartVPS(server *v2ray.Server) error {
	// 从全局配置获取 VeID 和 ApiKey
	bwgConfig := global.GVA_CONFIG.BWG
	if bwgConfig.VeID == "" || bwgConfig.ApiKey == "" {
		return fmt.Errorf("BWG 配置未设置，请在 config.yaml 中配置 bwg.veid 和 bwg.apiKey")
	}

	if err := (&StatService{}).PreCollectServerTraffic(server); err != nil {
		return err
	}

	// 构建请求 URL
	url := fmt.Sprintf("https://api.64clouds.com/v1/restart?veid=%s&api_key=%s", bwgConfig.VeID, bwgConfig.ApiKey)

	// 发送 HTTP GET 请求
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("重启失败 (HTTP %d): %s", resp.StatusCode, string(body))
	}

	// 解析 JSON 响应
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查是否有错误信息
	if errMsg, ok := result["error"].(string); ok && errMsg != "" {
		return fmt.Errorf("重启失败: %s", errMsg)
	}

	return nil
}
