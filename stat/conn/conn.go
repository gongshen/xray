package conn

import (
	"github.com/sirupsen/logrus"
	statsservice "github.com/xtls/xray-core/app/stats/command"
	"google.golang.org/grpc"
)

const TrafficTarget = "127.0.0.1:11111"

var closeFunc []func()
var StatServiceClient statsservice.StatsServiceClient

func InitConn() {
	closeFunc = make([]func(), 0)
	cc, err := grpc.Dial(TrafficTarget, grpc.WithInsecure())
	if err != nil {
		logrus.Errorln(err)
		return
	}
	closeFunc = append(closeFunc, func() {
		cc.Close()
	})
	StatServiceClient = statsservice.NewStatsServiceClient(cc)
}

func CloseConn() {
	for _, f := range closeFunc {
		f()
	}
}
