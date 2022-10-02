package main

import (
	"github.com/gongshen/xray/business"
	"github.com/gongshen/xray/conn"
)

func main() {
	if err := conn.InitDB("/usr/local/etc/xray.db"); err != nil {
		panic(err)
	}
	initLogger()
	srv := business.NewServer()
	srv.Start()
}
