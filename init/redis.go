package init

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"sephiroth-go/core"
)

func Redis() {
	redisCfg := core.Config.Redis
	var client redis.UniversalClient
	// 使用集群模式
	if redisCfg.UseCluster {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    redisCfg.ClusterAddrs,
			Password: redisCfg.Password,
		})
	} else {
		// 使用单例模式
		client = redis.NewClient(&redis.Options{
			Addr:     redisCfg.Addr,
			Password: redisCfg.Password,
			DB:       redisCfg.DB,
		})
	}
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		core.Log.Error("redis connect ping failed, err:", zap.Error(err))
		panic(err)
	} else {
		core.Log.Info("redis connect ping response:", zap.String("pong", pong))
		core.RedisClient = client
	}
}
