package main

import (
	"database/sql"
	"go.uber.org/zap"
	"sephiroth-go/core"
	"sephiroth-go/core/log"
	"sephiroth-go/initialize"
)

// //go:generate go env -w GO111MODULE=on
// //go:generate go env -w GOPROXY=https://goproxy.cn,direct
// //go:generate go mod tidy
// //go:generate go mod download

func main() {
	core.Vp = core.Viper()
	core.OtherInit()
	log.Log = core.Zap()
	zap.ReplaceGlobals(log.Log)
	core.Db = initialize.Database()
	initialize.Timer()
	initialize.DbList()

	if core.Db != nil {
		initialize.RegisterTables()
		db, _ := core.Db.DB()
		defer func(db *sql.DB) {
			_ = db.Close()
		}(db)
	}
	initialize.RunWindowsServer()
}
