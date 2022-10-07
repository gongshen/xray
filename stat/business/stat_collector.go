package business

import (
	"context"
	"github.com/gongshen/xray/stat/dao"
	"github.com/gongshen/xray/stat/models"
	"github.com/sirupsen/logrus"
	statsservice "github.com/xtls/xray-core/app/stats/command"
	"google.golang.org/grpc"
	"time"
)

type StatCollector struct{}

func (*StatCollector) Run() {
	cc, err := grpc.Dial(TrafficTarget, grpc.WithInsecure())
	if err != nil {
		logrus.Errorln(err)
		return
	}
	defer cc.Close()

	client := statsservice.NewStatsServiceClient(cc)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	request := &statsservice.QueryStatsRequest{
		Reset_: false,
	}
	resp, err := client.QueryStats(ctx, request)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	tagTrafficMap := map[string]*models.Traffic{}
	tags := make([]string, 0)
	for _, stat := range resp.GetStat() {
		matchs := trafficRegex.FindStringSubmatch(stat.Name)
		isUser := matchs[1] == "user"
		tag := matchs[2]
		//isDown := matchs[3] == "downlink"
		if tag == "api" || !isUser {
			continue
		}
		traffic, ok := tagTrafficMap[tag]
		if !ok {
			traffic = &models.Traffic{
				Tag: tag,
			}
			tagTrafficMap[tag] = traffic
			tags = append(tags, tag)
		}
		traffic.Used += stat.Value
	}
	if len(tags) <= 0 {
		return
	}
	curTraffics, err := dao.GetTrafficsByTags(tags)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	if len(curTraffics) != len(tagTrafficMap) {
		logrus.Errorln("traffics不匹配")
		return
	}
	for _, oldT := range curTraffics {
		// key一定存在
		newT := tagTrafficMap[oldT.Tag]
		// new used > old used
		if newT.Used >= oldT.Used {
			oldT.Used = newT.Used
		} else {
			oldT.Base += oldT.Used
			oldT.Used = newT.Used
		}
		if err = dao.UpdateTraffic(oldT); err != nil {
			logrus.Errorln(err)
			return
		}
	}
	return
}
