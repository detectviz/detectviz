package config

import (
	"os"
	"strconv"
	"sync"

	zap_adapter "github.com/detectviz/detectviz/internal/adapters/logger"
	"github.com/detectviz/detectviz/pkg/configtypes"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
	"go.uber.org/zap"
)

// defaultProvider provides a thread-safe configuration provider that reads from an internal map and falls back to environment variables.
// zh: defaultProvider 是一個具備內部快取與 thread-safe 的設定提供者，若未設定則回傳對應的環境變數值。
type defaultProvider struct {
	configMap       map[string]string
	notifierConfigs []configtypes.NotifierConfig
	cacheConfig     configtypes.CacheConfig // 快取模組的組態設定
	log             logger.Logger
	mu              sync.RWMutex
}

// NewDefaultProvider creates a new defaultProvider instance.
// zh: 建立並回傳一個預設的設定提供者實例。
func NewDefaultProvider() *defaultProvider {
	return &defaultProvider{
		configMap: make(map[string]string),
		notifierConfigs: []configtypes.NotifierConfig{
			{Name: "email", Type: "email", Target: "noreply@example.com", Enable: true},
			{Name: "slack", Type: "slack", Target: "https://hooks.slack.com/services/xxx", Enable: true},
			{Name: "webhook", Type: "webhook", Target: "https://example.com/webhook", Enable: false},
		},
		cacheConfig: configtypes.CacheConfig{
			Backend: "memory",
			Redis: configtypes.RedisConfig{
				Address:  "localhost:6379",
				Password: "",
				DB:       0,
			},
		},
		log: zap_adapter.NewZapLogger(zap.NewNop().Sugar()), // 預設使用 Zap Logger
	}
}

// Get retrieves the value associated with the given key.
// If the key is not present in configMap, it returns the value from the environment.
// zh: 取得指定 key 的設定值，若 configMap 未命中，則回傳對應環境變數值。
func (p *defaultProvider) Get(key string) string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if val, ok := p.configMap[key]; ok {
		return val
	}
	return os.Getenv(key)
}

// GetInt retrieves the configuration value as an integer.
// zh: 取得指定 key 的設定值，並轉換為 int（若轉換失敗則回傳 0）。
func (p *defaultProvider) GetInt(key string) int {
	val := p.Get(key)
	i, _ := strconv.Atoi(val)
	return i
}

// GetBool retrieves the configuration value as a boolean.
// zh: 取得指定 key 的設定值，並轉換為 bool（若轉換失敗則回傳 false）。
func (p *defaultProvider) GetBool(key string) bool {
	val := p.Get(key)
	b, _ := strconv.ParseBool(val)
	return b
}

// GetOrDefault returns the value associated with the key or a provided default value if not found.
// zh: 取得指定 key 的設定值，若為空字串則回傳提供的預設值。
func (p *defaultProvider) GetOrDefault(key, defaultVal string) string {
	val := p.Get(key)
	if val == "" {
		return defaultVal
	}
	return val
}

// GetNotifierConfigs returns the list of notifier configurations.
// zh: 回傳 notifier 設定的配置列表。
func (p *defaultProvider) GetNotifierConfigs() []configtypes.NotifierConfig {
	return p.notifierConfigs
}

// GetCacheConfig returns the cache configuration.
// zh: 回傳快取模組的組態設定。
func (p *defaultProvider) GetCacheConfig() configtypes.CacheConfig {
	return p.cacheConfig
}

// Reload is a no-op for defaultProvider.
// This method is a placeholder for config reload logic, and always returns nil.
// If reload is not supported, returns nil. (See interface documentation.)
// zh: 預留重新載入設定檔功能，目前尚未實作。若不支援，則回傳 nil。
func (p *defaultProvider) Reload() error {
	return nil
}

// Set assigns a key-value pair to the configMap.
// WARNING: This method is intended for testing purposes only. Do NOT use in production code!
// zh: 寫入設定鍵值對。⚠ 僅供測試用途，請勿於正式執行路徑中使用。
func (p *defaultProvider) Set(key, val string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.configMap[key] = val
}

// Logger returns the configured logger instance.
// zh: 回傳已配置的 logger 實例。
func (p *defaultProvider) Logger() logger.Logger {
	return p.log
}
