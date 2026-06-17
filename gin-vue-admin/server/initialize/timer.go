package initialize

import (
	"fmt"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/job"
	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/robfig/cron/v3"
)

const (
	defaultTrafficCollectInterval = "1h"
	defaultSysInfoCollectInterval = "5m"
)

func Timer() {
	if global.GVA_CONFIG.Timer.Start {
		for i := range global.GVA_CONFIG.Timer.Detail {
			go func(detail config.Detail) {
				var option []cron.Option
				if global.GVA_CONFIG.Timer.WithSeconds {
					option = append(option, cron.WithSeconds())
				}
				_, err := global.GVA_Timer.AddTaskByFunc("ClearDB", global.GVA_CONFIG.Timer.Spec, func() {
					err := utils.ClearTable(global.GVA_DB, detail.TableName, detail.CompareField, detail.Interval)
					if err != nil {
						fmt.Println("timer error:", err)
					}
				}, option...)
				if err != nil {
					fmt.Println("add timer error:", err)
				}
			}(global.GVA_CONFIG.Timer.Detail[i])
		}
	}
	// 添加收集流量定时任务
	go func() {
		time.Sleep(1 * time.Minute)
		location, _ := time.LoadLocation("Asia/Shanghai")
		trafficCollectSpec, ok := trafficCollectCronSpec(global.GVA_CONFIG.TrafficCollectInterval)
		if !ok {
			fmt.Printf("invalid traffic_collect_interval %q, use default %s\n", global.GVA_CONFIG.TrafficCollectInterval, defaultTrafficCollectInterval)
		}

		if _, err := global.GVA_Timer.AddTaskByJob("traffic_collect", trafficCollectSpec, job.CollectorJob{}, cron.WithLocation(location)); err != nil {
			fmt.Println("add timer error:", err)
		}
		sysInfoCollectSpec, ok := sysInfoCollectCronSpec(global.GVA_CONFIG.SysInfoCollectInterval)
		if !ok {
			fmt.Printf("invalid sysinfo_collect_interval %q, use default %s\n", global.GVA_CONFIG.SysInfoCollectInterval, defaultSysInfoCollectInterval)
		}
		go job.SysInfoCollectorJob{}.Run()
		if _, err := global.GVA_Timer.AddTaskByJob("sysinfo_collect", sysInfoCollectSpec, job.SysInfoCollectorJob{}, cron.WithLocation(location)); err != nil {
			fmt.Println("add timer error:", err)
		}
		if _, err := global.GVA_Timer.AddTaskByJob("calc_traffic_limit", "@every 10m", job.CalculateMonthlyTrafficLimitJob{}, cron.WithLocation(location)); err != nil {
			fmt.Println("add timer error:", err)
		}
		if _, err := global.GVA_Timer.AddTaskByJob("reset_traffic_limit", "@daily", job.ResetMonthlyTrafficLimitJob{}, cron.WithLocation(location)); err != nil {
			fmt.Println("add timer error:", err)
		}
		if _, err := global.GVA_Timer.AddTaskByJob("quota_reset", "@daily", job.QuotaResetJob{}, cron.WithLocation(location)); err != nil {
			fmt.Println("add timer error:", err)
		}
	}()

}

func trafficCollectCronSpec(interval string) (string, bool) {
	return durationCronSpec(interval, defaultTrafficCollectInterval)
}

func sysInfoCollectCronSpec(interval string) (string, bool) {
	return durationCronSpec(interval, defaultSysInfoCollectInterval)
}

func durationCronSpec(interval string, defaultInterval string) (string, bool) {
	interval = strings.TrimSpace(interval)
	if interval == "" {
		return "@every " + defaultInterval, true
	}
	duration, err := time.ParseDuration(interval)
	if err != nil || duration <= 0 {
		return "@every " + defaultInterval, false
	}
	return "@every " + interval, true
}
