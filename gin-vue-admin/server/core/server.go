package core

import (
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
	ListenAndServeTLS(certFile, keyFile string) error
}

func RunWindowsServer() {
	if global.GVA_CONFIG.System.UseMultipoint || global.GVA_CONFIG.System.UseRedis {
		// 初始化redis服务
		initialize.Redis()
	}

	// 从db加载jwt数据
	if global.GVA_DB != nil {
		system.LoadAll()
	}

	Router := initialize.Routers()

	address := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)
	s := initServer(address, Router)
	// 保证文本顺序输出
	time.Sleep(10 * time.Microsecond)
	global.GVA_LOG.Info("server run success on ", zap.String("address", address))

	fmt.Printf(`
	欢迎使用 gin-vue-admin
	默认前端文件运行地址:http://127.0.0.1%s
`, address)

	// 如果配置了证书文件，使用 HTTPS，否则使用 HTTP
	if global.GVA_CONFIG.System.CertFile != "" && global.GVA_CONFIG.System.KeyFile != "" {
		global.GVA_LOG.Error(s.ListenAndServeTLS(global.GVA_CONFIG.System.CertFile, global.GVA_CONFIG.System.KeyFile).Error())
	} else {
		global.GVA_LOG.Error(s.ListenAndServe().Error())
	}
}
