package global

import (
	"github.com/valyala/fasthttp"
	"net"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/utils/timer"
	"github.com/songzhibin97/gkit/cache/local_cache"

	"golang.org/x/sync/singleflight"

	"go.uber.org/zap"

	"github.com/flipped-aurora/gin-vue-admin/server/config"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	GVA_DB     *gorm.DB
	GVA_DBList map[string]*gorm.DB
	GVA_REDIS  *redis.Client
	GVA_CONFIG config.Server
	GVA_VP     *viper.Viper
	// GVA_LOG    *oplogging.Logger
	GVA_LOG                 *zap.Logger
	GVA_Timer               timer.Timer = timer.NewTimerTask()
	GVA_Concurrency_Control             = &singleflight.Group{}

	BlackCache local_cache.Cache
	lock       sync.RWMutex
	HTTP_CLI   = &fasthttp.Client{
		ReadTimeout:                   10 * time.Second,
		WriteTimeout:                  10 * time.Second,
		MaxIdleConnDuration:           30 * time.Minute,
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
		Dial: func(addr string) (net.Conn, error) {
			return (&fasthttp.TCPDialer{
				Concurrency:      10,
				DNSCacheDuration: time.Hour * 72,
			}).DialTimeout(addr, 10*time.Second)
		},
	}
)

// GetGlobalDBByDBName 通过名称获取db list中的db
func GetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	return GVA_DBList[dbname]
}

// MustGetGlobalDBByDBName 通过名称获取db 如果不存在则panic
func MustGetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	db, ok := GVA_DBList[dbname]
	if !ok || db == nil {
		panic("db no init")
	}
	return db
}
