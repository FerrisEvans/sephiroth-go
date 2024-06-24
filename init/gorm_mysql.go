package init

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sephiroth-go/core"
)

func GormMysql() *gorm.DB {
	config := core.Config.Mysql
	if config.Dbname == "" {
		return nil
	}
	mysqlConfig := mysql.Config{
		DSN:                       config.Dsn(), // DSN data source name
		DefaultStringSize:         191,          // string 类型字段的默认长度
		SkipInitializeWithVersion: false,        // 根据版本自动配置
	}

	if db, err := gorm.Open(mysql.New(mysqlConfig), Gorm.Config(config.Prefix, config.Singular)); err != nil {
		return nil
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE="+config.Engine)
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
		return db
	}
}
