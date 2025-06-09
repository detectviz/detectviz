package cachestore

import (
	memoryadapter "github.com/detectviz/detectviz/internal/adapters/cachestore/memory"
	redisadapter "github.com/detectviz/detectviz/internal/adapters/cachestore/redis"
	"github.com/detectviz/detectviz/pkg/configtypes"
	cachestoreiface "github.com/detectviz/detectviz/pkg/ifaces/cachestore"
	goredis "github.com/redis/go-redis/v9"
)

// RegisterCacheStore 根據設定註冊 CacheStore adapter。
// zh: 若 config 指定 redis，則註冊 redis adapter；否則預設使用記憶體快取。
func RegisterCacheStore(cfg configtypes.CacheConfig) cachestoreiface.CacheStore {
	if cfg.Backend == "redis" {
		client := goredis.NewClient(&goredis.Options{
			Addr:     cfg.Redis.Address,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
		return redisadapter.NewRedisCacheStore(client)
	}

	return memoryadapter.NewMemoryCacheStore()
}
