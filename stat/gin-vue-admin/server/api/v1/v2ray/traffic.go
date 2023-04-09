package v2ray

import (
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/v2ray"
	"github.com/valyala/fasthttp"
	"github.com/xtls/xray-core/app/stats/command"
	"go.uber.org/zap"
	"regexp"
	"time"
)

type Traffic struct {
	Category string // inbound outbound user
	Tag      string //
	Label    string // downlink uplink
	Value    int64
}

var trafficRegex = regexp.MustCompile("(inbound|outbound|user)>>>([^>]+)>>>traffic>>>(downlink|uplink)")

type CollectorJob struct{}

func (job CollectorJob) Run() {
	// 获取所有的服务器
	srvs, err := serverService.GetAllServer()
	if err != nil {
		global.GVA_LOG.Error("CollectorJob GetAllServer", zap.Error(err))
		return
	}
	now := time.Now()
	// 转换为东八时区
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println(err)
		return
	}
	zeroTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	for _, srv := range srvs {
		req, resp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
		req.SetRequestURI(fmt.Sprintf("http://%s:%d/stat/traffic", srv.Ip, global.GVA_CONFIG.STAT_PORT))
		if err = global.HTTP_CLI.Do(req, resp); err != nil {
			global.GVA_LOG.Error("CollectorJob Do", zap.Error(err))
			return
		}
		statsResp := new(command.QueryStatsResponse)
		if err = json.Unmarshal(resp.Body(), statsResp); err != nil {
			fmt.Println(string(resp.Body()))
			global.GVA_LOG.Error("CollectorJob Unmarshal", zap.Error(err))
			return
		}
		itemMap := make(map[string]*v2ray.Stat)
		for _, stat := range statsResp.Stat {
			matchs := trafficRegex.FindStringSubmatch(stat.Name)
			if len(matchs) != 4 {
				global.GVA_LOG.Error("CollectorJob FindStringSubmatch")
				return
			}
			item, ok := itemMap[matchs[2]]
			if !ok {
				item = &v2ray.Stat{
					Category:  matchs[1],
					Tag:       matchs[2],
					CreatedAt: zeroTime.Unix(),
				}
				itemMap[matchs[2]] = item
			}
			if matchs[3] == "downlink" {
				item.Down += uint64(stat.Value)
			} else {
				item.Up += uint64(stat.Value)
			}
		}
		if err = statService.StatsCollector(itemMap); err != nil {
			global.GVA_LOG.Error("CollectorJob StatsCollector")
			return
		}
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}
}
