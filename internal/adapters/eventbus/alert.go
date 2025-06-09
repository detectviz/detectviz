package eventbusadapter

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/event"
	"github.com/detectviz/detectviz/pkg/ifaces/eventbus"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// AlertLoggerHandler is a sample implementation of AlertEventHandler that logs the alert.
// zh: AlertLoggerHandler 是接收到告警觸發事件時記錄日誌的處理器實作範例。
type AlertLoggerHandler struct{}

// HandleAlertTriggered handles the alert triggered event by logging it.
// zh: 接收到告警事件後，透過 logger 模組輸出結構化日誌。
func (h *AlertLoggerHandler) HandleAlertTriggered(ctx context.Context, event event.AlertTriggeredEvent) error {
	log := logger.FromContext(ctx)
	log.WithFields(map[string]any{
		"alert_id":   event.AlertID,
		"rule_name":  event.RuleName,
		"level":      event.Level,
		"instance":   event.Instance,
		"metric":     event.Metric,
		"comparison": event.Comparison,
		"value":      event.Value,
		"threshold":  event.Threshold,
		"message":    event.Message,
	}).Info("[ALERT] " + event.Message)
	return nil
}

// alertHandlers 是所有已註冊的 AlertEventHandler 清單
var alertHandlers []eventbus.AlertEventHandler

// RegisterAlertHandler 用於讓 plugin 模組註冊自訂的告警事件處理器。
// zh: 提供 plugin 自動註冊機制，會將處理器加入全域列表。
func RegisterAlertHandler(handler eventbus.AlertEventHandler) {
	alertHandlers = append(alertHandlers, handler)
}

// LoadPluginAlertHandlers 回傳目前已註冊的 plugin AlertEventHandler 清單。
// zh: 在註冊器中載入 plugin 註冊的所有 alert handler。
func LoadPluginAlertHandlers() []eventbus.AlertEventHandler {
	return alertHandlers
}
