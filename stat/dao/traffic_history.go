package dao

import (
	"github.com/gongshen/xray/stat/conn"
	"github.com/gongshen/xray/stat/models"
)

func SaveHistories(histories ...*models.TrafficHistory) error {
	db := conn.GetDB()
	return db.Model(&models.TrafficHistory{}).CreateInBatches(histories, 10).Error
}
