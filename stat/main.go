package main

import (
	"flag"
	"github.com/gongshen/xray/stat/business"
	"github.com/gongshen/xray/stat/utils"
	"github.com/sirupsen/logrus"
)

var (
	domain string
)

func init() {
	flag.StringVar(&domain, "domain", "", "vless项目的域名")
}

func main() {
	flag.Parse()
	utils.SetDomain(domain)
	utils.SetIp()
	initLogger()
	logrus.Debugln("域名", domain)
	srv := business.NewServer()
	srv.Start()
}
