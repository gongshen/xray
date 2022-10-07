package main

import (
	"flag"
	"github.com/gongshen/xray/stat/business"
	"github.com/gongshen/xray/stat/utils"
	"github.com/sirupsen/logrus"
)

var (
	domain   string
	resetDay string
)

func init() {
	flag.StringVar(&domain, "domain", "", "vless项目的域名")
	flag.StringVar(&resetDay, "reset", "25", "流量重置日期，默认25号凌晨")
}

func main() {
	business.InitStat()
	flag.Parse()
	utils.SetDomain(domain)
	utils.SetIp()
	initLogger()
	logrus.Debugln("域名", domain)
	srv := business.NewServer(resetDay)
	srv.Start()
}
