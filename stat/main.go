package main

import (
	"github.com/gongshen/xray/stat/business"
	"github.com/gongshen/xray/stat/conn"
)

func main() {
	if err := conn.InitDB("/usr/local/etc/xray.db"); err != nil {
		panic(err)
	}
	initLogger()
	srv := business.NewServer()
	srv.Start()
}
