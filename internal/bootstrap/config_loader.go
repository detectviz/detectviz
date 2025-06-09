package bootstrap

import (
	"github.com/detectviz/detectviz/pkg/config"
	// ConfigProvider 介面
	// zh: ConfigProvider 介面
	configiface "github.com/detectviz/detectviz/pkg/ifaces/config"
)

// LoadConfig initializes and returns a ConfigProvider instance.
// zh: 初始化並回傳 ConfigProvider 設定供應器。
func LoadConfig() configiface.ConfigProvider {
	return config.NewDefaultProvider()
}
