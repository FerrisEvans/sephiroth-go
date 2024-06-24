package initialize

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	logs "log"
	"os"
	"sephiroth-go/config"
	"sephiroth-go/config/datasource"
	"sephiroth-go/core"
	"sephiroth-go/core/log"
	"time"
)

var Gorm = new(_gorm)

type _gorm struct{}

func Database() *gorm.DB {
	switch core.Config.System.DbType {
	case config.MySql:
		return GormMysql()
	case config.Postgres:
		return GormPgSql()
	default:
		return GormMysql()
	}
}

func RegisterTables() {
	db := core.Db
	err := db.AutoMigrate(
	// todo
	)
	if err != nil {
		log.Log.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}

	log.Log.Info("register table success")
}

// Config gorm 自定义配置
func (g *_gorm) Config(prefix string, singular bool) *gorm.Config {
	var general datasource.GeneralDB
	switch core.Config.System.DbType {
	case "mysql":
		general = core.Config.Mysql.GeneralDB
	case "pgsql":
		general = core.Config.Pgsql.GeneralDB
	default:
		general = core.Config.Mysql.GeneralDB
	}
	return &gorm.Config{
		Logger: logger.New(NewWriter(general, logs.New(os.Stdout, "\r\n", logs.LstdFlags)), logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      general.LogLevel(),
			Colorful:      true,
		}),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   prefix,
			SingularTable: singular,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}
}
