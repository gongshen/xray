package business

import (
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"net/http"
	"os"
)

func UpdateConfig(ctx *fasthttp.RequestCtx) {
	data := ctx.Request.Body()
	if err := os.WriteFile(XrayConfigFile, data, 0644); err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}
	//if err := Systemctl(SystemctlRestartOpt, ServiceNameXray); err != nil {
	//	ctx.Error(err.Error(), http.StatusBadRequest)
	//	return
	//}
	logrus.Debugln("收到配置文件 重启xray.")
	ctx.SuccessString("application/json", "OK")
	return
}
