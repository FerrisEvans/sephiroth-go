package main

import (
	"go.uber.org/zap"
	"sephiroth-go/core"
	"sephiroth-go/init"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

func main() {
	core.Vp = core.Viper()
	core.OtherInit()
	core.Log = core.Zap()
	zap.ReplaceGlobals(core.Log)
	core.Db = init.Database()
	init.Timer()
}
