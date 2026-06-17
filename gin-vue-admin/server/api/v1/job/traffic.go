package job

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/v2ray"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/service/v2ray_admin"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

var serviceGroup = service.ServiceGroupApp

type CollectorJob struct{}

func (job CollectorJob) Run() {
	srvs, err := serviceGroup.V2rayAdminServiceGroup.GetAllServer()
	if err != nil {
		global.GVA_LOG.Error("CollectorJob GetAllServer", zap.Error(err))
		return
	}
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println(err)
		return
	}
	now := time.Now().In(location)
	createdAt := now.Year()*10000 + int(now.Month())*100 + now.Day()
	for _, srv := range srvs {
		if err := serviceGroup.V2rayAdminServiceGroup.CollectServerTrafficWithRetry(srv, createdAt, 3); err != nil {
			global.GVA_LOG.Error("CollectorJob CollectServerTraffic", zap.Error(err), zap.String("ip", srv.Ip))
		}
	}
}

type SysInfoCollectorJob struct{}

func (job SysInfoCollectorJob) Run() {
	srvs, err := serviceGroup.V2rayAdminServiceGroup.GetAllServer()
	if err != nil {
		global.GVA_LOG.Error("SysInfoCollectorJob GetAllServer", zap.Error(err))
		return
	}
	for _, srv := range srvs {
		collectSysInfo(srv)
	}
}

func collectSysInfo(srv *v2ray.Server) {
	req, resp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(fmt.Sprintf("http://%s:%d/stat/sysinfo", srv.Ip, srv.GetStatPort()))
	if err := global.HTTP_CLI.Do(req, resp); err != nil {
		global.GVA_LOG.Error("SysInfoCollectorJob Do sysinfo", zap.Error(err), zap.String("ip", srv.Ip))
		return
	}

	body := resp.Body()
	if len(body) == 0 {
		global.GVA_LOG.Debug("SysInfoCollectorJob sysinfo: empty response", zap.String("ip", srv.Ip))
		return
	}

	sysInfo := new(v2ray_admin.SysInfo)
	if err := json.Unmarshal(body, sysInfo); err != nil {
		global.GVA_LOG.Error("SysInfoCollectorJob Unmarshal sysinfo", zap.Error(err), zap.String("body", string(body)), zap.String("ip", srv.Ip))
		return
	}

	if err := serviceGroup.V2rayAdminServiceGroup.UpdateServerSysInfo(srv.ID, sysInfo); err != nil {
		global.GVA_LOG.Error("SysInfoCollectorJob UpdateServerSysInfo", zap.Error(err))
		return
	}
}

type QuotaResetJob struct{}

func (job QuotaResetJob) Run() {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println(err)
		return
	}
	now := time.Now().In(location)
	srvs, err := serviceGroup.V2rayAdminServiceGroup.GetAllServer()
	if err != nil {
		global.GVA_LOG.Error("CollectorJob GetAllServer", zap.Error(err))
		return
	}
	zeroTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	for _, srv := range srvs {
		if v2ray_admin.IsTrafficResetDay(now, srv.ResetDate) {
			global.GVA_LOG.Info("server traffic quota reset", zap.String("server_ip", srv.Ip), zap.Uint64("used_quota", srv.TotalQuota))
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
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println(err)
		return
	}
	now := time.Now().In(location)
	global.GVA_LOG.Info("CollectorJob")

	bindings, err := serviceGroup.V2rayAdminServiceGroup.GetAllBindings()
	if err != nil {
		global.GVA_LOG.Error("CollectorJob GetAllBindings", zap.Error(err))
		return
	}

	needUpdateServer := make(map[int]*v2ray.Server)
	for _, binding := range bindings {
		if binding == nil || binding.UserID <= 0 || binding.ServerID <= 0 {
			continue
		}

		shouldLimit := v2ray_admin.ShouldLimitTraffic(0, binding.User.TrafficLimit)
		if binding.User.TrafficLimit > 0 {
			startCreatedAt := v2ray_admin.MonthlyTrafficLimitStartCreatedAt(now)
			traffic, err := serviceGroup.V2rayAdminServiceGroup.UserTrafficSince(startCreatedAt, strconv.Itoa(binding.UserID), binding.Server.Ip)
			if err != nil {
				global.GVA_LOG.Error("CollectorJob UserTrafficSince:", zap.Error(err))
				return
			}
			shouldLimit = v2ray_admin.ShouldLimitTraffic(traffic, binding.User.TrafficLimit)
		}

		if binding.IsLimited == shouldLimit {
			continue
		}
		if err = serviceGroup.V2rayAdminServiceGroup.SetBindingTrafficLimit(binding.ID, shouldLimit); err != nil {
			global.GVA_LOG.Error("CollectorJob SetBindingTrafficLimit:", zap.Error(err))
			return
		}
		needUpdateServer[binding.ServerID] = &binding.Server
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
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println(err)
		return
	}
	now := time.Now().In(location)
	if !v2ray_admin.IsMonthlyTrafficLimitResetDay(now) {
		return
	}
	srvs, err := serviceGroup.V2rayAdminServiceGroup.GetAllServer()
	if err != nil {
		global.GVA_LOG.Error("ResetMonthlyTrafficLimitJob GetAllServer:", zap.Error(err))
		return
	}
	for _, srv := range srvs {
		if err = serviceGroup.V2rayAdminServiceGroup.ResetTrafficLimitByServerID(srv.ID); err != nil {
			global.GVA_LOG.Error("ResetMonthlyTrafficLimitJob ResetTrafficLimitByServerID:", zap.Error(err))
			return
		}
		if err = serviceGroup.V2rayAdminServiceGroup.ReportBinding(srv); err != nil {
			global.GVA_LOG.Error("ResetMonthlyTrafficLimitJob ReportBinding:", zap.Error(err))
			return
		}
	}
}
