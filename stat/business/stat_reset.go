package business

import (
	"github.com/gongshen/xray/stat/dao"
	"github.com/gongshen/xray/stat/models"
	"github.com/sirupsen/logrus"
	"time"
)

type StatReset struct{}

func (*StatReset) Run() {
	traffics, err := dao.GetTraffics()
	if err != nil {
		logrus.Errorln(err)
	}
	now := time.Now()
	createDay := now.Format("2006-01-02")
	histories := make([]*models.TrafficHistory, 0, len(traffics))
	for _, traffic := range traffics {
		histories = append(histories, &models.TrafficHistory{
			Used:       traffic.Used + traffic.Base,
			Tag:        traffic.Tag,
			CreatedDay: createDay,
		})
	}
	if len(histories) > 0 {
		if err := dao.SaveHistories(histories...); err != nil {
			logrus.Errorln(err)
			return
		}
	}
	if err := dao.ResetTraffics(); err != nil {
		logrus.Errorln(err)
		return
	}
}
