package business

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gongshen/xray/stat/models"
	"github.com/sirupsen/logrus"
	statsservice "github.com/xtls/xray-core/app/stats/command"
	"github.com/xtls/xray-core/common/net"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

type GinServer struct {
	l net.Listener
}

func NewGinServer() *GinServer {
	return &GinServer{}
}

func (g *GinServer) Start() error {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		return err
	}
	g.l = listener

	engine := gin.New()
	srv := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: engine,
	}
	// 获取流量统计信息
	engine.GET("/stat/get", GetTraffic)
	engine.GET("/tag/new", NewTag)
	engine.GET("/tag/del", DelTag)
	engine.GET("/share/link", ShareLink)
	engine.GET("/share/qr", ShareQr)

	go srv.Serve(listener)
	return nil
}

func (g *GinServer) Close() {
	g.l.Close()
}

func ShareLink(c *gin.Context) {
	tag := c.Query("tag")
	if tag == "" {
		c.JSON(500, errors.New("参数错误"))
		return
	}

}

func ShareQr(c *gin.Context) {
	tag := c.Query("tag")
	if tag == "" {
		c.JSON(500, errors.New("参数错误"))
		return
	}

}

func GetTraffic(c *gin.Context) {
	cc, err := grpc.Dial(TrafficTarget, grpc.WithInsecure())
	if err != nil {
		c.JSON(500, err)
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
		c.JSON(500, err)
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
		logrus.Debugln("tag:%s,used:%d", tag, stat.Value)
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
		output.WriteString(fmt.Sprintf("用户:%s<br>使用流量:%s<br>", traffic.Tag, format(traffic.Used)))
		output.WriteString("-----------------------------------------------------------------------<br>")
	}

	c.Data(200, "text/html", output.Bytes())
}

func format(used int64) string {
	if used < GB {
		return fmt.Sprintf("%.2fMB", float64(used)/MB)
	} else {
		return fmt.Sprintf("%.2fGB", float64(used)/GB)
	}
}
