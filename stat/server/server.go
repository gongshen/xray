package server

import (
	"fmt"

	"github.com/gongshen/xray/stat/business"
	"github.com/gongshen/xray/stat/utils"
	"github.com/valyala/fasthttp"
)

func StartServer(port int) error {
	h := requestHandler
	s := &fasthttp.Server{
		Handler:            h,
		DisableKeepalive:   false,
		ReadBufferSize:     1024 * 4,
		WriteBufferSize:    1024 * 4,
		MaxRequestBodySize: 1024 * 1024 * 2,
	}
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Server listening on %s\n", addr)
	return s.ListenAndServe(addr)
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	remoteIP := ctx.RemoteIP().String()
	if remoteIP != utils.RemoteIp {
		ctx.Error("Forbidden", fasthttp.StatusForbidden)
		return
	}
	switch string(ctx.Path()) {
	case "/stat/traffic":
		business.CollectTraffic(ctx)
	case "/stat/traffic/collect":
		business.CollectTrafficToLocalStoreHandler(ctx)
	case "/stat/traffic/sync":
		business.SyncLocalTraffic(ctx)
	case "/stat/traffic/user-minute":
		business.AnalyzeUserTrafficHandler(ctx)
	case "/stat/sysinfo":
		business.CollectSysInfo(ctx)
	case "/conf/update":
		business.UpdateConfig(ctx)
	default:
		ctx.Error("Forbidden", fasthttp.StatusForbidden)
	}
}
