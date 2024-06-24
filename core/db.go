package core

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Db          *gorm.DB
	DbList      map[string]*gorm.DB
	RedisClient redis.UniversalClient
	//GVA_MONGO *qmgo.QmgoClient
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
		panic("db no initialize")
	}
	return db
}
