package v2ray

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/v2ray"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
    v2rayReq "github.com/flipped-aurora/gin-vue-admin/server/model/v2ray/request"
)

type StatService struct {
}

// CreateStat 创建Stat记录
// Author [piexlmax](https://github.com/piexlmax)
func (statService *StatService) CreateStat(stat *v2ray.Stat) (err error) {
	err = global.GVA_DB.Create(stat).Error
	return err
}

// DeleteStat 删除Stat记录
// Author [piexlmax](https://github.com/piexlmax)
func (statService *StatService)DeleteStat(stat v2ray.Stat) (err error) {
	err = global.GVA_DB.Delete(&stat).Error
	return err
}

// DeleteStatByIds 批量删除Stat记录
// Author [piexlmax](https://github.com/piexlmax)
func (statService *StatService)DeleteStatByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]v2ray.Stat{},"id in ?",ids.Ids).Error
	return err
}

// UpdateStat 更新Stat记录
// Author [piexlmax](https://github.com/piexlmax)
func (statService *StatService)UpdateStat(stat v2ray.Stat) (err error) {
	err = global.GVA_DB.Save(&stat).Error
	return err
}

// GetStat 根据id获取Stat记录
// Author [piexlmax](https://github.com/piexlmax)
func (statService *StatService)GetStat(id uint) (stat v2ray.Stat, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&stat).Error
	return
}

// GetStatInfoList 分页获取Stat记录
// Author [piexlmax](https://github.com/piexlmax)
func (statService *StatService)GetStatInfoList(info v2rayReq.StatSearch) (list []v2ray.Stat, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&v2ray.Stat{})
    var stats []v2ray.Stat
    // 如果有条件搜索 下方会自动创建搜索语句
    if info.StartCreatedAt !=nil && info.EndCreatedAt !=nil {
     db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
    }
    if info.Category != "" {
        db = db.Where("category = ?",info.Category)
    }
    if info.Tag != "" {
        db = db.Where("tag = ?",info.Tag)
    }
    if info.Created_day != nil {
        db = db.Where("created_day = ?",info.Created_day)
    }
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	err = db.Limit(limit).Offset(offset).Find(&stats).Error
	return  stats, total, err
}
