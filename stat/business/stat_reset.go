package business

import (
	"github.com/gongshen/xray/stat/dao"
	"github.com/sirupsen/logrus"
)

type StatReset struct{}

func (*StatReset) Run() {
	err := dao.ResetTraffics()
	if err != nil {
		logrus.Errorln(err)
		return
	}
}
