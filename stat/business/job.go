package business

import (
	"github.com/robfig/cron/v3"
	"time"
)

type JobServer struct {
	c *cron.Cron
}

func NewJobServer() *JobServer {
	return &JobServer{}
}

func (j *JobServer) Start() error {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return err
	}
	c := cron.New(cron.WithLocation(location), cron.WithSeconds())
	c.Start()
	go func() {
		time.Sleep(time.Second * 5)
		// 每 10 秒统计一次流量，首次启动延迟 5 秒，与重启 xray 的时间错开
		c.AddJob("@every 10s", NewTrafficService())
	}()
	return nil
}
func (j *JobServer) Close() {
	j.c.Stop()
}