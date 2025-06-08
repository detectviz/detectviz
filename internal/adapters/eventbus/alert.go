package eventbus

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/eventbus"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// AlertLoggerHandler is a sample implementation of AlertEventHandler that logs the alert.
// zh: AlertLoggerHandler 是接收到告警觸發事件時記錄日誌的處理器實作範例。
type AlertLoggerHandler struct{}

// HandleAlertTriggered handles the alert triggered event by logging it.
// zh: 接收到告警事件後，透過 logger 模組輸出結構化日誌。
func (h *AlertLoggerHandler) HandleAlertTriggered(ctx context.Context, event eventbus.AlertTriggeredEvent) error {
	log := logger.FromContext(ctx)
	log.WithFields(map[string]any{
		"rule_id":  event.RuleID,
		"severity": event.Severity,
		"message":  event.Message,
	}).Info("[ALERT] " + event.Message)
	return nil
}
