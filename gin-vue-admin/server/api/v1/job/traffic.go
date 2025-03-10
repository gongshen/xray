package job

import (
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/v2ray"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/valyala/fasthttp"
	"github.com/xtls/xray-core/app/stats/command"
	"go.uber.org/zap"
	"regexp"
	"strconv"
	"time"
)

var serviceGroup = service.ServiceGroupApp

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
	srvs, err := serviceGroup.V2rayAdminServiceGroup.GetAllServer()
	if err != nil {
		global.GVA_LOG.Error("CollectorJob GetAllServer", zap.Error(err))
		return
	}
	// 转换为东八时区
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println(err)
		return
	}
	now := time.Now().In(location)
	year := now.Year()
	month := now.Month()
	day := now.Day()
	createdAt := year*10000 + int(month)*100 + day
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
		// TODO 不知道总额度是：user+inbound+outbound，还是sum(user) = inbound+outbound
		var usedQuota uint64
		itemMap := make(map[string]*v2ray.Stat)
		for _, stat := range statsResp.Stat {
			if stat.Value <= 0 {
				continue
			}
			matchs := trafficRegex.FindStringSubmatch(stat.Name)
			if len(matchs) != 4 {
				global.GVA_LOG.Error("CollectorJob FindStringSubmatch")
				return
			}
			if matchs[1] != "user" {
				continue
			}
			item, ok := itemMap[matchs[2]]
			if !ok {
				item = &v2ray.Stat{
					Tag:       matchs[2],
					CreatedAt: createdAt,
					ServerIp:  srv.Ip,
				}
				itemMap[matchs[2]] = item
			}
			if matchs[3] == "downlink" {
				item.Down += uint64(stat.Value)
			} else {
				item.Up += uint64(stat.Value)
			}
			usedQuota += uint64(stat.Value)
		}
		if err = serviceGroup.V2rayAdminServiceGroup.StatsCollector(itemMap); err != nil {
			global.GVA_LOG.Error("CollectorJob StatsCollector")
			return
		}
		if err = serviceGroup.V2rayAdminServiceGroup.UpdateServerUsedQuota(srv.ID, usedQuota); err != nil {
			global.GVA_LOG.Error("CollectorJob UpdateServerUsedQuota")
			return
		}
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}
}

type QuotaResetJob struct{}

func (job QuotaResetJob) Run() {
	// 转换为东八时区
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println(err)
		return
	}
	now := time.Now().In(location)
	// 获取所有的服务器
	srvs, err := serviceGroup.V2rayAdminServiceGroup.GetAllServer()
	if err != nil {
		global.GVA_LOG.Error("CollectorJob GetAllServer", zap.Error(err))
		return
	}
	zeroTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	for _, srv := range srvs {
		if srv.ResetDate == now.Day() {
			global.GVA_LOG.Info("服务器流量重置", zap.String("服务器ip", srv.Ip), zap.Uint64("使用流量", srv.TotalQuota))
			if err = serviceGroup.V2rayAdminServiceGroup.SaveServerUsedQuotaLog(&v2ray.ServerQuotaLog{
				ServerID:  int(srv.ID),
				UsedQuota: srv.UsedQuota,
				CreatedAt: zeroTime.Unix(),
			}); err != nil {
				global.GVA_LOG.Error("CollectorJob SaveServerUsedQuotaLog", zap.Error(err))
				return
			}
			if err = serviceGroup.V2rayAdminServiceGroup.ResetServerUsedQuota(srv.ID); err != nil {
				global.GVA_LOG.Error("CollectorJob ResetServerUsedQuota", zap.Error(err))
				return
			}
		}
	}
}

type CalculateMonthlyTrafficLimitJob struct{}

func (job CalculateMonthlyTrafficLimitJob) Run() {
	// TODO 确保是每月一号
	// 转换为东八时区
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println(err)
		return
	}
	now := time.Now().In(location)
	global.GVA_LOG.Info("CollectorJob")
	// 获取所有的用户
	users, err := service.ServiceGroupApp.SystemServiceGroup.UserService.GetAllUser()
	if err != nil {
		global.GVA_LOG.Error("CollectorJob GetAllUser", zap.Error(err))
		return
	}
	var needLimitUser []uint
	for _, user := range users {
		if user.TrafficLimit < 0 {
			continue
		}
		if user.TrafficLimit == 0 {
			needLimitUser = append(needLimitUser, user.ID)
			continue
		}
		traffic, err := serviceGroup.V2rayAdminServiceGroup.MonthlyUserTraffic(now, strconv.Itoa(int(user.ID)))
		if err != nil {
			global.GVA_LOG.Error("CollectorJob MonthlyGetUserStat:", zap.Error(err))
			return
		}
		if traffic >= uint64(user.TrafficLimit*1024*1024*1024) {
			needLimitUser = append(needLimitUser, user.ID)
		}
	}
	needUpdateServer := make(map[int]*v2ray.Server)
	for _, user := range needLimitUser {
		bindings, err := serviceGroup.V2rayAdminServiceGroup.GetBindingByUserID(user)
		if err != nil {
			global.GVA_LOG.Error("CollectorJob GetBindingByUserID:", zap.Error(err))
			return
		}
		if err = serviceGroup.V2rayAdminServiceGroup.SetTrafficLimit(user); err != nil {
			global.GVA_LOG.Error("CollectorJob SetTrafficLimit:", zap.Error(err))
			return
		}
		for _, binding := range bindings {
			needUpdateServer[binding.ServerID] = &binding.Server
		}
	}
	for _, srv := range needUpdateServer {
		if err = serviceGroup.V2rayAdminServiceGroup.ReportBinding(srv); err != nil {
			global.GVA_LOG.Error("CollectorJob ReportBinding:", zap.Error(err))
			return
		}
	}
}

type ResetMonthlyTrafficLimitJob struct{}

func (job ResetMonthlyTrafficLimitJob) Run() {
	if err := serviceGroup.V2rayAdminServiceGroup.ResetTrafficLimit(); err != nil {
		global.GVA_LOG.Error("ResetMonthlyTrafficLimitJob ResetTrafficLimit:", zap.Error(err))
		return
	}
	srvs, err := serviceGroup.V2rayAdminServiceGroup.GetAllServer()
	if err != nil {
		global.GVA_LOG.Error("ResetMonthlyTrafficLimitJob GetAllServer:", zap.Error(err))
		return
	}
	for _, srv := range srvs {
		if err = serviceGroup.V2rayAdminServiceGroup.ReportBinding(srv); err != nil {
			global.GVA_LOG.Error("ResetMonthlyTrafficLimitJob ReportBinding:", zap.Error(err))
			return
		}
	}
}
