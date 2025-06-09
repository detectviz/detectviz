package alertlog

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/event"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// AlertPluginHandler is a sample plugin that handles alert events.
// zh: AlertPluginHandler 是一個範例插件，實作告警事件處理邏輯。
type AlertPluginHandler struct{}

// HandleAlertTriggered logs or processes the alert event.
// zh: 在此處理 AlertTriggeredEvent，例如自訂通知或紀錄。
func (h *AlertPluginHandler) HandleAlertTriggered(ctx context.Context, e event.AlertTriggeredEvent) error {
	log := logger.FromContext(ctx)
	log.WithFields(map[string]any{
		"alert_id":   e.AlertID,
		"rule_name":  e.RuleName,
		"level":      e.Level,
		"instance":   e.Instance,
		"metric":     e.Metric,
		"comparison": e.Comparison,
		"value":      e.Value,
		"threshold":  e.Threshold,
	}).Info("[ALERT] triggered")
	return nil
}
