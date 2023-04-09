package v2ray

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/v2ray"
	v2rayReq "github.com/flipped-aurora/gin-vue-admin/server/model/v2ray/request"
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
func (bindingService *BindingService) DeleteBinding(binding v2ray.Binding) (err error) {
	err = global.GVA_DB.Delete(&binding).Error
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
	err = global.GVA_DB.Where("id = ?", id).First(&binding).Error
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
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Preload("User").Preload("Server").Find(&bindings).Error
	return bindings, total, err
}
