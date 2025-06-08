package config

import (
	"github.com/detectviz/detectviz/pkg/configtypes"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// ConfigProvider defines the interface for retrieving configuration values.
// zh: ConfigProvider 定義取得設定值的抽象介面，可支援從環境變數、檔案或遠端服務載入。
type ConfigProvider interface {
	// Get returns the string value for a given key.
	// zh: 根據指定 key 取得對應的字串設定值。
	Get(key string) string

	// GetInt returns the integer value for a given key.
	// zh: 根據指定 key 取得整數型別的設定值。
	GetInt(key string) int

	// GetBool returns the boolean value for a given key.
	// zh: 根據指定 key 取得布林型別的設定值。
	GetBool(key string) bool

	// GetOrDefault returns the value for a given key, or returns the provided default if not found.
	// zh: 根據 key 取得設定值，若無對應值則回傳預設值。
	GetOrDefault(key, defaultVal string) string

	// GetNotifierConfigs returns the list of notifier configurations.
	// zh: 回傳 notifier 設定的配置列表。
	GetNotifierConfigs() []configtypes.NotifierConfig

	// Logger returns the configured logger instance.
	// zh: 回傳已配置的 logger 實例。
	Logger() logger.Logger

	// Reload refreshes the underlying configuration source, if supported.
	// zh: 重新載入設定來源內容（若支援），用於動態更新設定。
	Reload() error
}
