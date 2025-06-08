package configtypes

// CacheConfig defines the configuration for the cache module.
// zh: 支援記憶體與 Redis 快取，透過 backend 參數選擇。
type CacheConfig struct {
	Backend string      `json:"backend" yaml:"backend"` // The cache backend to use ("memory" or "redis") // 使用的快取後端（memory 或 redis）
	Redis   RedisConfig `json:"redis" yaml:"redis"`     // Redis configuration // Redis 設定
}

// RedisConfig defines the configuration for Redis cache.
// zh: 僅在 backend 設為 redis 時會使用。
type RedisConfig struct {
	Address  string `json:"address" yaml:"address"`   // Redis connection address // Redis 連線位址
	Password string `json:"password" yaml:"password"` // Redis password (can be empty) // Redis 密碼（可留空）
	DB       int    `json:"db" yaml:"db"`             // Redis DB number (default 0) // Redis 使用的 DB 編號（預設 0）
}
