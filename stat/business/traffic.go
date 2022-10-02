package business

import (
	"context"
	"fmt"
	"github.com/gongshen/xray/stat/conn"
	"github.com/gongshen/xray/stat/dao"
	"github.com/gongshen/xray/stat/models"
	"github.com/sirupsen/logrus"
	statsservice "github.com/xtls/xray-core/app/stats/command"
	"google.golang.org/grpc"
	"regexp"
	"time"
)

const defaultLimit = 100 * 1024 * 1024 * 1024 // 默认100G

var trafficRegex = regexp.MustCompile("(inbound|outbound|user)>>>([^>]+)>>>traffic>>>(downlink|uplink)")

func DisableInvalidTraffic() (int64, error) {
	db := conn.GetDB()
	result := db.Model(models.Traffic{}).
		Where("up + down >= limit and enable = ?", true).
		Update("enable", false)
	err := result.Error
	count := result.RowsAffected
	return count, err
}

type TrafficService struct{}

func NewTrafficService() *TrafficService {
	return &TrafficService{}
}

func (*TrafficService) Run() {
	c, err := grpc.Dial(fmt.Sprintf("127.0.0.1:11111"), grpc.WithInsecure())
	if err != nil {
		logrus.Println(err)
		return
	}
	defer c.Close()

	client := statsservice.NewStatsServiceClient(c)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	request := &statsservice.QueryStatsRequest{
		Reset_: false,
	}
	resp, err := client.QueryStats(ctx, request)
	if err != nil {
		logrus.Println(err)
		return
	}
	tagTrafficMap := map[string]*models.Traffic{}
	traffics := make([]*models.Traffic, 0)
	for _, stat := range resp.GetStat() {
		logrus.Debugf("打印xray信息: name:%s,value:%d", stat.Name, stat.Value)
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
			traffics = append(traffics, traffic)
		}
		traffic.Used += stat.Value
	}
	if err = dao.UpdateTraffic(traffics); err != nil {
		logrus.Println(err)
		return
	}
}
