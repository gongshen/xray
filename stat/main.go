package main

import (
	"github.com/gongshen/xray/stat/business"
)

func main() {
	initLogger()
	srv := business.NewServer()
	srv.Start()
}
