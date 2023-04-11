package business

import (
	"context"
	"encoding/json"
	"github.com/gongshen/xray/stat/conn"
	"github.com/valyala/fasthttp"
	statsservice "github.com/xtls/xray-core/app/stats/command"
	"net/http"
	"time"
)

func CollectTraffic(reqCtx *fasthttp.RequestCtx) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	request := &statsservice.QueryStatsRequest{
		Reset_: true,
	}
	resp, err := conn.StatServiceClient.QueryStats(ctx, request)
	if err != nil {
		reqCtx.Error(err.Error(), http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(resp)
	if err != nil {
		reqCtx.Error(err.Error(), http.StatusBadRequest)
		return
	}
	reqCtx.Success("application/json", data)
	return
}
