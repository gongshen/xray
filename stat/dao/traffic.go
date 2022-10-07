package dao

import (
	"github.com/gongshen/xray/stat/conn"
	"github.com/gongshen/xray/stat/models"
)

// UpdateTraffic 用户的traffic更新
func UpdateTraffic(traffic *models.Traffic) (err error) {
	db := conn.GetDB()
	db = db.Model(models.Traffic{})
	cond := map[string]interface{}{
		"used": traffic.Used,
		"base": traffic.Base,
	}
	err = db.Where("tag = ?", traffic.Tag).
		UpdateColumns(cond).Error
	if err != nil {
		return
	}
	return
}

// NewTraffic 根据tag新建traffic
// 先检查该traffic是否存在
func NewTraffic(tag string) error {
	db := conn.GetDB()
	affected := db.Model(models.Traffic{}).Where("tag = ?", tag).UpdateColumn("enable", true).RowsAffected
	if affected <= 0 {
		return db.Model(models.Traffic{}).Save(&models.Traffic{
			Tag:    tag,
			Enable: true,
		}).Error
	}
	return nil
}

// DelTraffic 软删除traffic
func DelTraffic(tag string) error {
	db := conn.GetDB()
	return db.Model(models.Traffic{}).Where("tag = ?", tag).UpdateColumn("enable", false).Error
}

func GetEnableTraffics() ([]*models.Traffic, error) {
	db := conn.GetDB()
	var ans []*models.Traffic
	err := db.Model(models.Traffic{}).Where("enable = ?", true).Find(&ans).Error
	if err != nil && !conn.IsNotFound(err) {
		return nil, err
	}
	return ans, nil
}

func GetTrafficsByTags(tags []string) ([]*models.Traffic, error) {
	db := conn.GetDB()
	var ans []*models.Traffic
	err := db.Model(models.Traffic{}).Where("tag in (?)", tags).Find(&ans).Error
	if err != nil && !conn.IsNotFound(err) {
		return nil, err
	}
	return ans, nil
}
