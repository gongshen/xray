package dao

import (
	"github.com/gongshen/xray/conn"
	"github.com/gongshen/xray/models"
	"gorm.io/gorm"
)

func SaveTraffic(traffic *models.Traffic) error {
	db := conn.GetDB()
	return db.Model(models.Traffic{}).Save(traffic).Error
}

func UpdateTraffic(traffics []*models.Traffic) (err error) {
	db := conn.GetDB()
	db = db.Model(models.Traffic{})
	for _, traffic := range traffics {
		err = db.Debug().Where("tag = ?", traffic.Tag).
			UpdateColumn("used", gorm.Expr("used + ?", traffic.Used)).Error
		if err != nil {
			return
		}
	}
	return
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
