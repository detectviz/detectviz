package eventbus

import (
	// 強制載入 plugin 模組以觸發其 init() 註冊事件處理器
	adapterlogger "github.com/detectviz/detectviz/internal/adapters/logger"
	"github.com/detectviz/detectviz/pkg/ifaces/event"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

var (
	registeredAlertHandlers  []event.AlertEventHandler
	registeredHostHandlers   []event.HostEventHandler
	registeredMetricHandlers []event.MetricEventHandler
	registeredTaskHandlers   []event.TaskEventHandler
	defaultLogger            logger.Logger = adapterlogger.NewNopLogger()
)

// RegisterPluginAlertHandler allows plugin modules to register alert handlers.
// zh: 外部模組註冊自訂 Alert 處理器。
func RegisterPluginAlertHandler(h event.AlertEventHandler) {
	registeredAlertHandlers = append(registeredAlertHandlers, h)
}

// RegisterPluginHostHandler allows plugin modules to register host handlers.
// zh: 外部模組註冊自訂 Host 處理器。
func RegisterPluginHostHandler(h event.HostEventHandler) {
	registeredHostHandlers = append(registeredHostHandlers, h)
}

// RegisterPluginMetricHandler allows plugin modules to register metric handlers.
// zh: 外部模組註冊自訂 Metric 處理器。
func RegisterPluginMetricHandler(h event.MetricEventHandler) {
	registeredMetricHandlers = append(registeredMetricHandlers, h)
}

// RegisterPluginTaskHandler allows plugin modules to register task handlers.
// zh: 外部模組註冊自訂 Task 處理器。
func RegisterPluginTaskHandler(h event.TaskEventHandler) {
	registeredTaskHandlers = append(registeredTaskHandlers, h)
}

// LoadPluginAlertHandlers retrieves all plugin-registered alert handlers.
// zh: 載入所有插件註冊的 Alert 處理器。
func LoadPluginAlertHandlers() []event.AlertEventHandler {
	return registeredAlertHandlers
}

// LoadPluginHostHandlers retrieves all plugin-registered host handlers.
// zh: 載入所有插件註冊的 Host 處理器。
func LoadPluginHostHandlers() []event.HostEventHandler {
	return registeredHostHandlers
}

// LoadPluginMetricHandlers retrieves all plugin-registered metric handlers.
// zh: 載入所有插件註冊的 Metric 處理器。
func LoadPluginMetricHandlers() []event.MetricEventHandler {
	return registeredMetricHandlers
}

// LoadPluginTaskHandlers retrieves all plugin-registered task handlers.
// zh: 載入所有插件註冊的 Task 處理器。
func LoadPluginTaskHandlers() []event.TaskEventHandler {
	return registeredTaskHandlers
}

// ExplorePlugins 探索並初始化所有 event plugin 模組。
// zh: 匯入所有 plugin，以觸發各自 init 函式註冊對應事件處理器。
func ExplorePlugins() {
	// 模組透過上方 import _ 初始化，故此函式本身不需邏輯
}

// OverrideDefaultLogger 設定 plugin 使用的預設 logger。
// zh: 測試時可使用此方法注入 log 攔截器。
func OverrideDefaultLogger(log logger.Logger) {
	defaultLogger = log
}

// GetDefaultLogger 回傳 plugin 使用的 logger。
// zh: 提供 plugin 在 init 時注入的預設 logger 實例。
func GetDefaultLogger() logger.Logger {
	return defaultLogger
}
