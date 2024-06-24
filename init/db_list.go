package init

import (
	"gorm.io/gorm"
	"sephiroth-go/config/datasource"
	"sephiroth-go/core"
)

const sys = "system"

func DbList() {
	dbMap := make(map[string]*gorm.DB)
	for _, info := range core.Config.DbList {
		if info.Disable {
			continue
		}
		switch info.Type {
		case "mysql":
			dbMap[info.AliasName] = GormMysqlByConfig(datasource.Mysql{
				GeneralDB: info.GeneralDB,
			})
		case "pgsql":
			dbMap[info.AliasName] = GormPgSqlByConfig(datasource.Pgsql{
				GeneralDB: info.GeneralDB,
			})
		default:
			continue
		}
	}
	// 做特殊判断,是否有迁移
	// 适配低版本迁移多数据库版本
	if sysDB, ok := dbMap[sys]; ok {
		core.Db = sysDB
	}
	core.DbList = dbMap
}
