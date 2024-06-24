package init

import (
	"github.com/songzhibin97/gkit/cache/local_cache"
	"sephiroth-go/core"
	"sephiroth-go/util"
)

func OtherInit() {
	dr, err := util.ParseDuration(core.Config.Jwt.ExpiresTime)
	if err != nil {
		panic(err)
	}
	_, err = util.ParseDuration(core.Config.Jwt.BufferTime)
	if err != nil {
		panic(err)
	}

	core.BlackCache = local_cache.NewCache(
		local_cache.SetDefaultExpire(dr),
	)
}
