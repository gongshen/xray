package business

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gongshen/xray/stat/models"
	statsservice "github.com/xtls/xray-core/app/stats/command"
	"google.golang.org/grpc"
	"time"
)

func GetStat(c *gin.Context) {
	reset := c.Query("reset")
	cc, err := grpc.Dial(TrafficTarget, grpc.WithInsecure())
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	defer cc.Close()

	client := statsservice.NewStatsServiceClient(cc)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	request := &statsservice.QueryStatsRequest{
		Reset_: reset == "",
	}
	resp, err := client.QueryStats(ctx, request)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	tagTrafficMap := map[string]*models.Traffic{}
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
		}
		traffic.Used += stat.Value
	}
	var output bytes.Buffer
	output.WriteString("<html><body><p>")
	for _, traffic := range tagTrafficMap {
		output.WriteString(fmt.Sprintf(`<font color="red">用户：</font>%s<br><font color="green">已用流量：</font>%s<br>`, traffic.Tag, format(traffic.Used)))
		output.WriteString("---------------------------------------<br>")
	}
	output.WriteString("</p></body></html>")
	c.Data(200, "text/html; charset=utf-8", output.Bytes())
}

func format(used int64) string {
	if used < GB {
		return fmt.Sprintf("%.2fMB", float64(used)/MB)
	} else {
		return fmt.Sprintf("%.2fGB", float64(used)/GB)
	}
}
