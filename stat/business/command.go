package business

import (
	"github.com/sirupsen/logrus"
	"os/exec"
)

const (
	SystemctlRestartOpt = "restart"
	SystemctlStartOpt   = "start"
	SystemctlStopOpt    = "stop"
)

const (
	ServiceNameXray  = "xray"
	ServiceNameNginx = "nginx"
	ServiceNameStat  = "stat"
)

func Systemctl(opt string, svcName string) error {
	cmd := exec.Command("systemctl", opt, svcName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	logrus.Debugf("%s %s:%s\n", svcName, opt, string(output))
	return nil
}
