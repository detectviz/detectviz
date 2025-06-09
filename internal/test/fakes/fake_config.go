// Package fakes 提供測試用的假物件。
// zh: 用於測試的設定假實作，支援 key-value 模擬與錯誤注入。
package fakes

import (
	"errors"
	"strconv"

	"github.com/detectviz/detectviz/pkg/configtypes"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// FakeConfig implements the ConfigProvider interface for testing purposes.
// zh: 用於測試的設定假實作，支援 key-value 模擬與錯誤注入。
type FakeConfig struct {
	// Values 儲存 key-value 字串對應，模擬設定內容。
	Values map[string]string
	// CacheCfg 模擬快取設定。
	CacheCfg configtypes.CacheConfig
	// NotifierCfgs 模擬通知器設定清單。
	NotifierCfgs []configtypes.NotifierConfig
	// ReloadCalled 標記 Reload 是否被呼叫過。
	ReloadCalled bool
	// LoggerInstance 模擬 logger 實例。
	LoggerInstance logger.Logger
	// ReloadShouldFail 控制 Reload 是否回傳錯誤。
	ReloadShouldFail bool
}

// Get returns the string value for a given key.
// zh: 取得指定 key 的字串值，若不存在則回傳空字串。
func (f *FakeConfig) Get(key string) string {
	return f.Values[key]
}

// GetOrDefault returns the value or default if not set.
// zh: 取得指定 key 的字串值，若不存在則回傳預設值。
func (f *FakeConfig) GetOrDefault(key, defaultVal string) string {
	val, ok := f.Values[key]
	if !ok {
		return defaultVal
	}
	return val
}

// GetInt returns the int value for a given key.
// zh: 取得指定 key 的整數值，若格式錯誤則回傳 0。
func (f *FakeConfig) GetInt(key string) int {
	v, _ := strconv.Atoi(f.Values[key])
	return v
}

// GetBool returns the bool value for a given key.
// zh: 取得指定 key 的布林值，若格式錯誤則回傳 false。
func (f *FakeConfig) GetBool(key string) bool {
	v, _ := strconv.ParseBool(f.Values[key])
	return v
}

// GetCacheConfig returns the cache config struct.
// zh: 取得快取設定結構。
func (f *FakeConfig) GetCacheConfig() configtypes.CacheConfig {
	return f.CacheCfg
}

// GetNotifierConfigs returns notifier configuration list.
// zh: 取得通知器設定清單。
func (f *FakeConfig) GetNotifierConfigs() []configtypes.NotifierConfig {
	return f.NotifierCfgs
}

// Logger returns the logger instance.
// zh: 取得 logger 實例。
func (f *FakeConfig) Logger() logger.Logger {
	return f.LoggerInstance
}

// Reload sets ReloadCalled and optionally returns error.
// zh: 重新載入設定，會標記 ReloadCalled，並可模擬錯誤。
func (f *FakeConfig) Reload() error {
	f.ReloadCalled = true
	if f.ReloadShouldFail {
		return errors.New("simulated reload failure")
	}
	return nil
}
