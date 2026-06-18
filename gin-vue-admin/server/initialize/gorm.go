package initialize

import (
	"os"

	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"

	"github.com/flipped-aurora/gin-vue-admin/server/model/v2ray"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Gorm() *gorm.DB {
	switch global.GVA_CONFIG.System.DbType {
	case "mysql":
		return GormMysql()
	case "pgsql":
		return GormPgSql()
	case "oracle":
		return GormOracle()
	case "mssql":
		return GormMssql()
	default:
		return GormMysql()
	}
}

func RegisterTables() {
	db := global.GVA_DB
	err := db.AutoMigrate(
		// 系统模块表
		system.SysApi{},
		system.SysUser{},
		system.SysBaseMenu{},
		system.JwtBlacklist{},
		system.SysAuthority{},
		system.SysDictionary{},
		system.SysOperationRecord{},
		system.SysDictionaryDetail{},
		system.SysBaseMenuParameter{},
		system.SysBaseMenuBtn{},
		system.SysAuthorityBtn{},
		system.SysChatGptOption{},
		adapter.CasbinRule{},

		example.ExaFile{},
		example.ExaCustomer{},
		example.ExaFileChunk{},
		example.ExaFileUploadAndDownload{},
		v2ray.Stat{},
		v2ray.StatSnapshot{},
		v2ray.StatSyncCursor{},
		v2ray.TrafficAnomaly{},
		v2ray.Server{},
		v2ray.Binding{},
		v2ray.ServerQuotaLog{},
	)
	if err != nil {
		global.GVA_LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}
	ensureManagedApi(db, "/v2ray_admin/server/analyzeUserTraffic", "POST", "分析用户流量明细")
	ensureManagedApi(db, "/v2ray_admin/server/classifyTrafficTargets", "POST", "访问目标分类")
	global.GVA_LOG.Info("register table success")
}

func ensureManagedApi(db *gorm.DB, path string, method string, description string) {
	apiAttrs := system.SysApi{
		ApiGroup:    "v2ray_admin",
		Method:      method,
		Path:        path,
		Description: description,
	}
	var api system.SysApi
	if err := db.Where("path = ? AND method = ?", path, method).Attrs(apiAttrs).FirstOrCreate(&api).Error; err != nil {
		global.GVA_LOG.Warn("ensure managed api failed", zap.Error(err), zap.String("path", path), zap.String("method", method))
	}

	ruleAttrs := adapter.CasbinRule{Ptype: "p", V0: "9527", V1: path, V2: method}
	var rule adapter.CasbinRule
	if err := db.Where(&ruleAttrs).Attrs(ruleAttrs).FirstOrCreate(&rule).Error; err != nil {
		global.GVA_LOG.Warn("ensure managed api casbin rule failed", zap.Error(err), zap.String("path", path), zap.String("method", method))
	}
}
