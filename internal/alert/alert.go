package alert

import (
	ifconfig "github.com/detectviz/detectviz/pkg/ifaces/config"
)

var cfg ifconfig.ConfigProvider

func Init(c ifconfig.ConfigProvider) {
	cfg = c
}

// 範例：使用設定值啟用或停用告警模組
func IsAlertEnabled() bool {
	return cfg.GetBool("alert.enabled")
}
