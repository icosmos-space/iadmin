package global

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/qmgo"

	"github.com/icosmos-space/iadmin/server/utils/timer"
	"github.com/songzhibin97/gkit/cache/local_cache"

	"golang.org/x/sync/singleflight"

	"go.uber.org/zap"

	"github.com/icosmos-space/iadmin/server/config"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	IADMIN_DB     *gorm.DB
	IADMIN_DBList map[string]*gorm.DB
	IADMIN_REDIS  redis.UniversalClient
	GVA_REDISList map[string]redis.UniversalClient
	GVA_MONGO     *qmgo.QmgoClient
	IADMIN_CONFIG config.Server
	IADMIN_VP     *viper.Viper
	// IADMIN_LOG    *oplogging.Logger
	IADMIN_LOG              *zap.Logger
	IADMIN_Timer            timer.Timer = timer.NewTimerTask()
	GVA_Concurrency_Control             = &singleflight.Group{}
	IADMIN_ROUTERS          gin.RoutesInfo
	IADMIN_ACTIVE_DBNAME    *string
	BlackCache              local_cache.Cache
	lock                    sync.RWMutex
)

// GetGlobalDBByDBName 通过名称获取db list中的db
func GetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	return IADMIN_DBList[dbname]
}

// MustGetGlobalDBByDBName 通过名称获取db 如果不存在则panic
func MustGetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	db, ok := IADMIN_DBList[dbname]
	if !ok || db == nil {
		panic("db no init")
	}
	return db
}

func GetRedis(name string) redis.UniversalClient {
	redis, ok := GVA_REDISList[name]
	if !ok || redis == nil {
		panic(fmt.Sprintf("redis `%s` no init", name))
	}
	return redis
}
