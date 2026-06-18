package server

import (
	"encoding/json"
	"net"
	"path/filepath"
	"testing"

	"github.com/gongshen/xray/stat/business"
	"github.com/gongshen/xray/stat/utils"
	"github.com/valyala/fasthttp"
)

func TestAccountedRequestHandlerRecordsRequestAndResponseBytes(t *testing.T) {
	store, err := business.OpenTrafficStore(filepath.Join(t.TempDir(), "stat.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	previousStore := business.LocalStore
	business.LocalStore = store
	defer func() {
		business.LocalStore = previousStore
	}()

	previousRemoteIP := utils.RemoteIp
	utils.RemoteIp = "127.0.0.1"
	defer func() {
		utils.RemoteIp = previousRemoteIP
	}()

	var req fasthttp.Request
	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetRequestURI("/unknown")
	req.Header.SetHost("stat.example")
	req.Header.SetUserAgent("stat-test")
	req.SetBodyString("request-body")

	var ctx fasthttp.RequestCtx
	ctx.Init(&req, &net.TCPAddr{IP: net.ParseIP(utils.RemoteIp), Port: 12345}, nil)

	accountedRequestHandler(&ctx)

	events, err := store.ListEventsAfter(0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fatalf("events length = %d, want 1", len(events))
	}
	if events[0].Tag != "1" {
		t.Fatalf("tag = %q, want 1", events[0].Tag)
	}
	if events[0].Up <= uint64(len("request-body")) {
		t.Fatalf("up = %d, want greater than request body length", events[0].Up)
	}
	if events[0].Down <= 0 {
		t.Fatalf("down = %d, want positive response bytes", events[0].Down)
	}
}

func TestRequestHandlerAllowsLoopbackForTrafficEventIngest(t *testing.T) {
	store, err := business.OpenTrafficStore(filepath.Join(t.TempDir(), "stat.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	previousStore := business.LocalStore
	business.LocalStore = store
	defer func() {
		business.LocalStore = previousStore
	}()

	previousRemoteIP := utils.RemoteIp
	utils.RemoteIp = "173.242.124.21"
	defer func() {
		utils.RemoteIp = previousRemoteIP
	}()

	body, err := json.Marshal(business.TrafficEventIngestRequest{
		Events: []business.LocalTrafficEvent{{Tag: "1", Down: 10, Up: 5}},
	})
	if err != nil {
		t.Fatal(err)
	}

	var req fasthttp.Request
	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetRequestURI("/stat/traffic/event")
	req.SetBody(body)

	var ctx fasthttp.RequestCtx
	ctx.Init(&req, &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 12345}, nil)

	requestHandler(&ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusOK {
		t.Fatalf("status = %d, body = %s", ctx.Response.StatusCode(), ctx.Response.Body())
	}
}
