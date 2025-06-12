package config

import (
	"github.com/detectviz/detectviz/pkg/configtypes"
	"github.com/spf13/viper"
)

// DefaultConfigProvider 為預設的 ConfigProvider 實作。
type DefaultConfigProvider struct {
	v *viper.Viper
}

// NewDefaultConfigProvider 建立 Viper-based 的 ConfigProvider。
func NewDefaultConfigProvider(v *viper.Viper) *DefaultConfigProvider {
	return &DefaultConfigProvider{v: v}
}

func (p *DefaultConfigProvider) GetLoggerConfig() configtypes.LoggerConfig {
	var c configtypes.LoggerConfig
	_ = p.v.UnmarshalKey("logging", &c)
	return c
}

func (p *DefaultConfigProvider) GetSchedulerConfig() configtypes.SchedulerConfig {
	var c configtypes.SchedulerConfig
	_ = p.v.UnmarshalKey("scheduler", &c)
	return c
}

func (p *DefaultConfigProvider) GetAlertConfig() configtypes.AlertConfig {
	var c configtypes.AlertConfig
	_ = p.v.UnmarshalKey("alert", &c)
	return c
}

func (p *DefaultConfigProvider) GetBusConfig() configtypes.BusConfig {
	var c configtypes.BusConfig
	_ = p.v.UnmarshalKey("eventbus", &c)
	return c
}

func (p *DefaultConfigProvider) GetEncryptionConfig() configtypes.EncryptionConfig {
	var c configtypes.EncryptionConfig
	_ = p.v.UnmarshalKey("security.encryption", &c)
	return c
}

func (p *DefaultConfigProvider) GetWebConfig() configtypes.WebConfig {
	var c configtypes.WebConfig
	_ = p.v.UnmarshalKey("web", &c)
	return c
}

func (p *DefaultConfigProvider) GetPluginConfig(key string, target interface{}) error {
	return p.v.UnmarshalKey("plugin."+key, target)
}
