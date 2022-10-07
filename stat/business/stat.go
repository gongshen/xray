package business

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gongshen/xray/stat/conn"
	"github.com/gongshen/xray/stat/dao"
	"github.com/gongshen/xray/stat/models"
	"github.com/sirupsen/logrus"
	statsservice "github.com/xtls/xray-core/app/stats/command"
	"google.golang.org/grpc"
	"sort"
	"time"
)

func InitStat() {
	if err := conn.InitDB(XrayDBPath); err != nil {
		panic(err)
	}
}

func GetStat(c *gin.Context) {
	traffics, err := dao.GetEnableTraffics()
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	sort.Slice(traffics, func(i, j int) bool {
		return traffics[i].Base+traffics[i].Used > traffics[j].Base+traffics[j].Used
	})
	var output bytes.Buffer
	output.WriteString("<html><body>")
	output.WriteString("---------------------------------------<br>")
	for _, traffic := range traffics {
		output.WriteString(fmt.Sprintf(`<font color="red">用户：</font>%s<br><font color="green">已用流量：</font>%s<br>`, traffic.Tag, format(traffic.Used+traffic.Base)))
		output.WriteString("---------------------------------------<br>")
	}
	output.WriteString("</body></html>")
	c.Data(200, "text/html; charset=utf-8", output.Bytes())
}

func format(used int64) string {
	if used < GB {
		return fmt.Sprintf("%.2fMB", float64(used)/MB)
	} else {
		return fmt.Sprintf("%.2fGB", float64(used)/GB)
	}
}

type Stat struct{}

func (*Stat) Run() {
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
