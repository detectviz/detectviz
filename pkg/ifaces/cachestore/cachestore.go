package cachestore

// CacheStore defines the abstract interface for caching used in Detectviz.
// It supports basic key-value operations with TTL (time-to-live) and can be implemented using in-memory or Redis-based stores.
// zh: CacheStore 定義 Detectviz 中的快取抽象介面，支援基本的鍵值操作與 TTL（存活時間），可對接 memory 或 Redis 實作。
type CacheStore interface {
	// Get returns the cached value for the given key, or an error if not found.
	// zh: 根據 key 取得對應快取內容，若 key 不存在應回傳錯誤。
	Get(key string) (string, error)

	// Set sets the cache value for a given key with a TTL in seconds.
	// If ttlSeconds is 0, the value never expires.
	// zh: 設定快取內容與存活時間，ttlSeconds 為 0 時表示永不過期。
	Set(key string, val string, ttlSeconds int) error

	// Has returns true if the given key exists in cache.
	// zh: 檢查 key 是否存在於快取中，非強一致性。
	Has(key string) bool

	// Delete removes the value associated with the given key.
	// zh: 從快取中移除指定的 key。
	Delete(key string) error

	// Keys returns all cache keys that start with the given prefix.
	// Optional method for grouping or bulk invalidation use cases.
	// zh: 回傳所有以指定 prefix 開頭的 key，常用於群組清除。
	Keys(prefix string) ([]string, error)
}
