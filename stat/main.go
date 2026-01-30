package main

import (
	"flag"

	"github.com/gongshen/xray/stat/conn"
	"github.com/gongshen/xray/stat/server"
	"github.com/gongshen/xray/stat/utils"
	"github.com/sirupsen/logrus"
)

var (
	//	domain   string
	//	resetDay string
	level string
	port  int
)

func init() {
	//flag.StringVar(&domain, "domain", "", "vless项目的域名")
	//flag.StringVar(&resetDay, "reset", "25", "流量重置日期，默认25号凌晨")
	flag.StringVar(&level, "level", "info", "日志级别")
	flag.IntVar(&port, "port", 56611, "监听端口")
}

func main() {
	flag.Parse()
	utils.SetRemoteIp()
	utils.SetIp()
	lv, _ := logrus.ParseLevel(level)
	logrus.SetLevel(lv)
	conn.InitConn()
	defer conn.CloseConn()
	if err := server.StartServer(port); err != nil {
		logrus.Println(err)
		return
	}
}
