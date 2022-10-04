package main

import (
	"flag"
	"fmt"
	"github.com/gongshen/xray/stat/business"
	"github.com/gongshen/xray/stat/utils"
	"github.com/sirupsen/logrus"
	"net/url"
)

var (
	domain string
)

func init() {
	flag.StringVar(&domain, "domain", "", "vless项目的域名")
}

func main() {
	fmt.Println(url.PathEscape("vless://$UUID@$DOMAIN:$PORT?security=xtls&flow=$FLOW#XTLS_wulabing-$DOMAIN"))
	return
	flag.Parse()
	logrus.Debugln("域名", domain)
	utils.SetDomain(domain)
	initLogger()
	srv := business.NewServer()
	srv.Start()
}
