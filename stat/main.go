package main

import (
	"flag"
	"time"

	"github.com/gongshen/xray/stat/business"
	"github.com/gongshen/xray/stat/conn"
	"github.com/gongshen/xray/stat/server"
	"github.com/gongshen/xray/stat/utils"
	"github.com/sirupsen/logrus"
)

var (
	level           string
	port            int
	trafficDBPath   string
	collectInterval time.Duration
)

func init() {
	flag.StringVar(&level, "level", "info", "log level")
	flag.IntVar(&port, "port", 56611, "listen port")
	flag.StringVar(&trafficDBPath, "traffic-db", "/var/lib/xray-stat/stat.db", "traffic sqlite db path")
	flag.DurationVar(&collectInterval, "collect-interval", 5*time.Second, "traffic collect interval")
}

func main() {
	flag.Parse()
	utils.SetRemoteIp()
	utils.SetIp()
	lv, _ := logrus.ParseLevel(level)
	logrus.SetLevel(lv)

	conn.InitConn()
	defer conn.CloseConn()

	store, err := business.OpenTrafficStore(trafficDBPath)
	if err != nil {
		logrus.Println(err)
		return
	}
	business.LocalStore = store
	defer store.Close()

	stopCollector := make(chan struct{})
	go business.StartLocalTrafficCollector(store, collectInterval, stopCollector)
	defer close(stopCollector)

	if err := server.StartServer(port); err != nil {
		logrus.Println(err)
		return
	}
}
