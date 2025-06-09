package cachestoreadapter

import (
	"errors"
	"fmt"
	"sort"

	memoryadapter "github.com/detectviz/detectviz/internal/adapters/cachestore/memory"
	redisadapter "github.com/detectviz/detectviz/internal/adapters/cachestore/redis"
	"github.com/detectviz/detectviz/pkg/ifaces/cachestore"
	goredis "github.com/redis/go-redis/v9"
)

// DefaultBackend 為預設的快取實作。
// zh: 若無指定 backend，預設使用 memory。
const DefaultBackend = "memory"

// ErrUnknownBackend 表示快取 backend 未註冊。
// zh: 用於回傳尚未註冊的 backend 錯誤。
var ErrUnknownBackend = errors.New("unknown cache backend")

var registry = map[string]func() cachestore.CacheStore{
	"memory": func() cachestore.CacheStore {
		return memoryadapter.NewMemoryCacheStore()
	},
}

// Register registers a new cache backend.
// zh: 註冊一個新的快取後端。
func Register(name string, fn func() cachestore.CacheStore) {
	registry[name] = fn
}

// Get returns a registered cache backend by name.
// zh: 根據名稱取得對應的快取實作，若不存在則回傳錯誤。
func Get(name string) (cachestore.CacheStore, error) {
	if fn, ok := registry[name]; ok {
		return fn(), nil
	}
	return nil, fmt.Errorf("%w: %s", ErrUnknownBackend, name)
}

// GetDefault returns the default cache backend.
// zh: 取得預設的快取後端實作。
func GetDefault() cachestore.CacheStore {
	fn := registry[DefaultBackend]
	return fn()
}

// WithRedisClient registers and returns a Redis-based cache store.
// zh: 註冊並建立 Redis 快取實作，用於注入 redis client。
func WithRedisClient(client *goredis.Client) cachestore.CacheStore {
	Register("redis", func() cachestore.CacheStore {
		return redisadapter.NewRedisCacheStore(client)
	})
	return redisadapter.NewRedisCacheStore(client)
}

// List returns the names of all registered cache backends.
// zh: 回傳所有已註冊的快取後端名稱清單（依字母排序）。
func List() []string {
	var names []string
	for name := range registry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
