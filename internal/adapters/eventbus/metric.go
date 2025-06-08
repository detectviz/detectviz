package eventbus

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/event"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// MetricLoggerHandler is a sample implementation of MetricEventHandler that logs overflow events.
// zh: MetricLoggerHandler 是接收到指標溢出事件時記錄日誌的處理器實作範例。
type MetricLoggerHandler struct{}

// HandleMetricOverflow handles MetricOverflowEvent by logging the event with structured fields.
// zh: 接收到指標溢出事件後，使用結構化欄位輸出警告日誌
func (h *MetricLoggerHandler) HandleMetricOverflow(ctx context.Context, event event.MetricOverflowEvent) error {
	log := logger.FromContext(ctx)
	log.WithFields(map[string]any{
		"metric":    event.MetricName,
		"value":     event.Value,
		"threshold": event.Threshold,
		"instance":  event.Instance,
		"reason":    event.Reason,
	}).Warn("[METRIC] overflow detected")
	return nil
}

// metricHandlers 是所有已註冊的 MetricEventHandler 清單
var metricHandlers []event.MetricEventHandler

// RegisterMetricHandler 用於讓 plugin 模組註冊自訂的指標事件處理器。
// zh: 提供 plugin 自動註冊機制，會將處理器加入全域列表。
func RegisterMetricHandler(handler event.MetricEventHandler) {
	metricHandlers = append(metricHandlers, handler)
}

// LoadPluginMetricHandlers 回傳目前已註冊的 plugin MetricEventHandler 清單。
// zh: 在註冊器中載入 plugin 註冊的所有 metric handler。
func LoadPluginMetricHandlers() []event.MetricEventHandler {
	return metricHandlers
}
