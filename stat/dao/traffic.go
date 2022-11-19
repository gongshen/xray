package dao

import (
	"github.com/gongshen/xray/stat/conn"
	"github.com/gongshen/xray/stat/models"
)

// UpdateTraffic 用户的traffic更新
func UpdateTraffic(traffic *models.Traffic) error {
	return conn.GetDB().Model(models.Traffic{}).Exec("update traffics set used = used + ? where tag = ?", traffic.Used, traffic.Tag).Error
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

func GetTraffics() ([]*models.Traffic, error) {
	db := conn.GetDB()
	var ans []*models.Traffic
	err := db.Model(models.Traffic{}).Find(&ans).Error
	if err != nil && !conn.IsNotFound(err) {
		return nil, err
	}
	return ans, nil
}

func ResetTraffics() error {
	db := conn.GetDB()
	// 保存到history表
	if err := db.Model(models.Traffic{}).Where("enable = ?", 1).Update("base", 0).Error; err != nil {
		return err
	}
	if err := db.Model(models.Traffic{}).Where("enable = ?", 1).Update("used", 0).Error; err != nil {
		return err
	}
	return nil
}
