package core

import (
	"github.com/songzhibin97/gkit/cache/local_cache"
	"sephiroth-go/util"
)

func OtherInit() {
	dr, err := util.ParseDuration(Config.Jwt.ExpiresTime)
	if err != nil {
		panic(err)
	}
	_, err = util.ParseDuration(Config.Jwt.BufferTime)
	if err != nil {
		panic(err)
	}

	BlackCache = local_cache.NewCache(
		local_cache.SetDefaultExpire(dr),
	)
}
