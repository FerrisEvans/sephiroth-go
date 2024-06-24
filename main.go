package main

import (
	"sephiroth-go/core"
	"sephiroth-go/init"
)

func main() {
	core.Vp = core.Viper()
	core.OtherInit()
	core.Log = core.Zap()
	core.Db = init.Database()
}
