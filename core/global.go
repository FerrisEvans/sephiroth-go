package core

import (
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/spf13/viper"
	"golang.org/x/sync/singleflight"
	"sephiroth-go/config"
	"sephiroth-go/core/timer"
	"sync"
)

var (
	Config             config.Server
	Vp                 *viper.Viper
	Timer              = timer.NewTimerTask()
	ConcurrencyControl = &singleflight.Group{}
	BlackCache         local_cache.Cache
	lock               sync.RWMutex
)
