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
	level                  string
	port                   int
	trafficDBPath          string
	collectInterval        time.Duration
	trafficRetentionMonths int
	statAPITrafficTag      string
	logCleanupDir          string
	logRetentionMonths     int
	xrayLogDir             string
	xrayLogRetentionMonths int
)

func init() {
	flag.StringVar(&level, "level", "info", "log level")
	flag.IntVar(&port, "port", 56611, "listen port")
	flag.StringVar(&trafficDBPath, "traffic-db", "/var/lib/xray-stat/stat.db", "traffic sqlite db path")
	flag.DurationVar(&collectInterval, "collect-interval", 10*time.Second, "traffic collect interval")
	flag.IntVar(&trafficRetentionMonths, "traffic-retention-months", 12, "traffic sqlite event retention months")
	flag.StringVar(&statAPITrafficTag, "stat-api-traffic-tag", business.DefaultStatAPITrafficTag, "traffic tag used for stat api request/response bytes")
	flag.StringVar(&logCleanupDir, "log-clean-dir", "/root/log", "xray-admin date directory cleanup root")
	flag.IntVar(&logRetentionMonths, "log-retention-months", 12, "xray-admin date directory retention months")
	flag.StringVar(&xrayLogDir, "xray-log-dir", "/var/log/xray", "xray log directory")
	flag.IntVar(&xrayLogRetentionMonths, "xray-log-retention-months", 12, "xray rotated log retention months")
}

func main() {
	flag.Parse()
	business.SetStatAPITrafficTag(statAPITrafficTag)
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

	stopTrafficEventCleaner := make(chan struct{})
	go business.StartTrafficEventCleaner(store, trafficRetentionMonths, 24*time.Hour, stopTrafficEventCleaner)
	defer close(stopTrafficEventCleaner)

	stopLogCleaner := make(chan struct{})
	go business.StartLogDirectoryCleaner(logCleanupDir, logRetentionMonths, 24*time.Hour, stopLogCleaner)
	defer close(stopLogCleaner)

	stopXrayLogCleaner := make(chan struct{})
	go business.StartXrayLogFileCleaner(xrayLogDir, xrayLogRetentionMonths, 24*time.Hour, stopXrayLogCleaner)
	defer close(stopXrayLogCleaner)

	if err := server.StartServer(port); err != nil {
		logrus.Println(err)
		return
	}
}
