package initialize

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/v2ray_admin"
	"github.com/robfig/cron/v3"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
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

		if _, err := global.GVA_Timer.AddTaskByJob("traffic_collect", "@every 10s", v2ray_admin.CollectorJob{}, cron.WithLocation(location)); err != nil {
			fmt.Println("add timer error:", err)
		}
		if _, err := global.GVA_Timer.AddTaskByJob("quota_reset", "@daily", v2ray_admin.QuotaResetJob{}, cron.WithLocation(location)); err != nil {
			fmt.Println("add timer error:", err)
		}
	}()

}
