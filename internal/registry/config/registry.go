package config

import (
	"github.com/detectviz/detectviz/pkg/config"
	configiface "github.com/detectviz/detectviz/pkg/ifaces/config"
)

// RegisterConfigProvider 初始化並回傳預設的 ConfigProvider。
// zh: 用於統一初始化設定模組，供其他模組依賴注入。
func RegisterConfigProvider() configiface.ConfigProvider {
	return config.NewDefaultProvider()
}
