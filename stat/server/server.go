package server

import (
	"fmt"
	"github.com/gongshen/xray/stat/business"
	"github.com/gongshen/xray/stat/utils"
	"github.com/valyala/fasthttp"
)

func StartServer() error {
	h := requestHandler
	s := &fasthttp.Server{
		Handler:            h,
		DisableKeepalive:   false, // 开启长连接支持
		ReadBufferSize:     1024 * 4,
		WriteBufferSize:    1024 * 4,
		MaxRequestBodySize: 1024 * 1024 * 2, // 最大请求体大小为 2MB
	}
	fmt.Println("Server listening on :56611")
	return s.ListenAndServe(":56611")
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
	case "/conf/update":
		business.UpdateConfig(ctx)
	default:
		ctx.Error("Forbidden", fasthttp.StatusForbidden)
	}
}
