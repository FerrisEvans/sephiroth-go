package init

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"sephiroth-go/config/datasource"
	"sephiroth-go/constant"
	"sephiroth-go/core"
	"time"
)

var Gorm = new(_gorm)

type _gorm struct{}

func Database() *gorm.DB {
	switch core.Config.System.DbType {
	case constant.MySql:
		return GormMysql()
	case constant.Postgres:
		return GormPgSql()
	default:
		return GormMysql()
	}
}

// Config gorm 自定义配置
// Author [SliverHorn](https://github.com/SliverHorn)
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
		Logger: logger.New(NewWriter(general, log.New(os.Stdout, "\r\n", log.LstdFlags)), logger.Config{
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
