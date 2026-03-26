package initialize

import (
	"context"
	"sync"

	"github.com/icosmos-space/iadmin/server/config"
	"github.com/icosmos-space/iadmin/server/global"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	redisOnce     sync.Once
	redisListOnce sync.Once
	redisInitMu   sync.Mutex
)

func initRedisClient(redisCfg config.Redis) (redis.UniversalClient, error) {
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
		global.IADMIN_LOG.Error("redis connect ping failed, err:", zap.String("name", redisCfg.Name), zap.Error(err))
		return nil, err
	}

	global.IADMIN_LOG.Info("redis connect ping response:", zap.String("name", redisCfg.Name), zap.String("pong", pong))
	return client, nil
}

func Redis() {
	redisInitMu.Lock()
	defer redisInitMu.Unlock()

	redisOnce.Do(func() {
		redisClient, err := initRedisClient(global.IADMIN_CONFIG.Redis)
		if err != nil {
			panic(err)
		}
		global.IADMIN_REDIS = redisClient
	})
}

func RedisList() {
	redisInitMu.Lock()
	defer redisInitMu.Unlock()

	redisListOnce.Do(func() {
		redisMap := make(map[string]redis.UniversalClient)

		for _, redisCfg := range global.IADMIN_CONFIG.RedisList {
			client, err := initRedisClient(redisCfg)
			if err != nil {
				panic(err)
			}
			redisMap[redisCfg.Name] = client
		}

		global.GVA_REDISList = redisMap
	})
}

// resetRedisSingleton 重置 Redis 单例状态，仅供测试辅助入口调用。
func resetRedisSingleton() {
	redisInitMu.Lock()
	defer redisInitMu.Unlock()

	if global.IADMIN_REDIS != nil {
		_ = global.IADMIN_REDIS.Close()
		global.IADMIN_REDIS = nil
	}

	for name, client := range global.GVA_REDISList {
		if client != nil {
			_ = client.Close()
		}
		delete(global.GVA_REDISList, name)
	}
	global.GVA_REDISList = nil

	redisOnce = sync.Once{}
	redisListOnce = sync.Once{}
}
