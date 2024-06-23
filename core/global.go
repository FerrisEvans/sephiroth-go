package core

import (
	"github.com/redis/go-redis/v9"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"sephiroth-go/config"
	"sephiroth-go/util"
	"sync"
)

var (
	Db          *gorm.DB
	DbList      map[string]*gorm.DB
	RedisClient redis.UniversalClient
	// //GVA_MONGO *qmgo.QmgoClient
	Config             config.Server
	Vp                 *viper.Viper
	Log                *zap.Logger
	Timer              util.Timer = util.NewTimerTask()
	ConcurrencyControl            = &singleflight.Group{}
	BlackCache         local_cache.Cache
	lock               sync.RWMutex
)

// GetGlobalDbByDBName 通过名称获取db list中的db
func GetGlobalDbByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	return DbList[dbname]
}

// MustGetGlobalDbByDBName 通过名称获取db 如果不存在则panic
func MustGetGlobalDbByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	db, ok := DbList[dbname]
	if !ok || db == nil {
		panic("db no init")
	}
	return db
}
