package business

import (
	"github.com/sirupsen/logrus"
	"os/exec"
)

func XrayRestart() error {
	cmd := exec.Command("systemctl", "restart", "xray")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	logrus.Debugln("xray restart:", string(output))
	return nil
}
