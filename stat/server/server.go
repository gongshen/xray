package server

import (
	"fmt"
	"net"
	"strconv"

	"github.com/gongshen/xray/stat/business"
	"github.com/gongshen/xray/stat/utils"
	"github.com/valyala/fasthttp"
)

func StartServer(port int) error {
	h := accountedRequestHandler
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

func accountedRequestHandler(ctx *fasthttp.RequestCtx) {
	up := estimateHTTPRequestBytes(&ctx.Request)
	requestHandler(ctx)
	down := estimateHTTPResponseBytes(&ctx.Response)
	business.RecordStatAPIUsage(up, down)
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	remoteIP := ctx.RemoteIP().String()
	if !isAllowedRemoteIP(remoteIP) {
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
	case "/stat/traffic/event":
		business.IngestTrafficEvents(ctx)
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

func isAllowedRemoteIP(remoteIP string) bool {
	if remoteIP == utils.RemoteIp {
		return true
	}
	ip := net.ParseIP(remoteIP)
	return ip != nil && ip.IsLoopback()
}

func estimateHTTPRequestBytes(req *fasthttp.Request) uint64 {
	total := len(req.Header.Method()) + len(req.RequestURI()) + len(" HTTP/1.1\r\n")
	total += estimateRequestHeaderBytes(&req.Header)
	total += len("\r\n")
	total += len(req.Body())
	return uint64(total)
}

func estimateHTTPResponseBytes(resp *fasthttp.Response) uint64 {
	statusCode := resp.StatusCode()
	total := len("HTTP/1.1 ") + len(strconv.Itoa(statusCode)) + len(" ") + len(fasthttp.StatusMessage(statusCode)) + len("\r\n")
	total += estimateResponseHeaderBytes(&resp.Header)
	total += len("\r\n")
	total += len(resp.Body())
	return uint64(total)
}

func estimateRequestHeaderBytes(header *fasthttp.RequestHeader) int {
	total := 0
	header.VisitAll(func(key, value []byte) {
		total += len(key) + len(": ") + len(value) + len("\r\n")
	})
	return total
}

func estimateResponseHeaderBytes(header *fasthttp.ResponseHeader) int {
	total := 0
	header.VisitAll(func(key, value []byte) {
		total += len(key) + len(": ") + len(value) + len("\r\n")
	})
	return total
}
