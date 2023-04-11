package v2ray_admin

import (
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/v2ray"
	v2rayReq "github.com/flipped-aurora/gin-vue-admin/server/model/v2ray/request"
	"github.com/valyala/fasthttp"
	"strconv"
)

type BindingService struct {
}

// CreateBinding 创建Binding记录
// Author [piexlmax](https://github.com/piexlmax)
func (bindingService *BindingService) CreateBinding(binding *v2ray.Binding) (err error) {
	err = global.GVA_DB.Create(binding).Error
	return err
}

// DeleteBinding 删除Binding记录
// Author [piexlmax](https://github.com/piexlmax)
func (bindingService *BindingService) DeleteBinding(binding *v2ray.Binding) (err error) {
	err = global.GVA_DB.Delete(binding).Error
	return err
}

// DeleteBindingByIds 批量删除Binding记录
// Author [piexlmax](https://github.com/piexlmax)
func (bindingService *BindingService) DeleteBindingByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]v2ray.Binding{}, "id in ?", ids.Ids).Error
	return err
}

// UpdateBinding 更新Binding记录
// Author [piexlmax](https://github.com/piexlmax)
func (bindingService *BindingService) UpdateBinding(binding v2ray.Binding) (err error) {
	err = global.GVA_DB.Save(&binding).Error
	return err
}

// GetBinding 根据id获取Binding记录
// Author [piexlmax](https://github.com/piexlmax)
func (bindingService *BindingService) GetBinding(id uint) (binding v2ray.Binding, err error) {
	err = global.GVA_DB.Where("id = ?", id).Preload("User").Preload("Server").First(&binding).Error
	return
}

// GetBindingInfoList 分页获取Binding记录
// Author [piexlmax](https://github.com/piexlmax)
func (bindingService *BindingService) GetBindingInfoList(info v2rayReq.BindingSearch) (list []v2ray.Binding, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&v2ray.Binding{})
	var bindings []v2ray.Binding
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.ServerID > 0 {
		db = db.Where("server_id = ?", info.ServerID)
	}
	if info.UserID > 0 {
		db = db.Where("user_id = ?", info.UserID)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Preload("User").Preload("Server").Find(&bindings).Error
	return bindings, total, err
}

func (bindingService *BindingService) CountBindingByServerID(ids ...int) (count int64, err error) {
	err = global.GVA_DB.Model(&v2ray.Binding{}).Where("server_id in (?)", ids).Count(&count).Error
	return
}

func (bindingService *BindingService) GetBindingsByServerID(ids ...int) (bindings []*v2ray.Binding, err error) {
	err = global.GVA_DB.Model(&v2ray.Binding{}).Where("server_id in (?)", ids).Preload("User").
		Preload("Server").Find(&bindings).Error
	return
}

func (bindingService *BindingService) CountBindingByUserID(ids ...int) (count int64, err error) {
	err = global.GVA_DB.Model(&v2ray.Binding{}).Where("user_id in (?)", ids).Count(&count).Error
	return
}

func (bindingService *BindingService) GetAllBindings() (bindings []*v2ray.Binding, err error) {
	err = global.GVA_DB.Preload("User").Preload("Server").Find(&bindings).Error
	return
}

// ReportBinding 上报绑定
func (bindingService *BindingService) ReportBinding(srv *v2ray.Server) error {
	// 获取所有该服务器的绑定
	bindings, err := bindingService.GetBindingsByServerID(int(srv.ID))
	if err != nil {
		return err
	}
	settingsCli := make([]*v2ray.XrayConfigSettingsClient, 0, len(bindings))
	for _, binding := range bindings {
		settingsCli = append(settingsCli, &v2ray.XrayConfigSettingsClient{
			Id:      binding.User.UUID.String(),
			Level:   binding.Level,
			AlterId: binding.AlterID,
			Email:   strconv.Itoa(binding.UserID),
		})
	}
	config := v2ray.NewXRayConfig(srv.Port)
	config.InboundConfigs[0].Settings.Clients = settingsCli
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	srv.Config = data
	req, resp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	req.Header.SetMethod("POST")
	req.SetBody(data)
	req.SetRequestURI(fmt.Sprintf("http://%s:%d/conf/update", srv.Ip, global.GVA_CONFIG.STAT_PORT))
	if err = global.HTTP_CLI.Do(req, nil); err != nil {
		return err
	}
	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)
	return nil
}
